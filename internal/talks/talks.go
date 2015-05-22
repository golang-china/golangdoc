// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package talks

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/present"
)

var validJSONPFunc = regexp.MustCompile(`(?i)^[a-z_][a-z0-9_.]*$`)

// Config specifies Server configuration values.
type Config struct {
	RootFS       vfs.FileSystem
	ContentPath  string // Relative or absolute location of article files and related content.
	TemplatePath string // Relative or absolute location of template files.

	BaseURL  string // Absolute base URL (for permalinks; no trailing slash).
	BasePath string // Base URL path relative to server root (no trailing slash).
	GodocURL string // The base URL of godoc (for menu bar; no trailing slash).
	Hostname string // Server host name, used for rendering ATOM feeds.

	PlayEnabled bool
}

// Doc represents an article adorned with presentation data.
type Doc struct {
	*present.Doc
	Permalink string        // Canonical URL for this document.
	Path      string        // Path relative to server root (including base).
	HTML      template.HTML // rendered article

	Related      []*Doc
	Newer, Older *Doc
}

// Server implements an http.Handler that serves blog articles.
type Server struct {
	cfg      Config
	docs     []*Doc
	tags     []string
	docPaths map[string]*Doc // key is path without BasePath.
	docTags  map[string][]*Doc
	template struct {
		action, article, dir, slides, doc *template.Template
	}

	// contentTemplate maps the presentable file extensions to the
	// template to be executed.
	contentTemplate map[string]*template.Template

	content http.Handler
}

// NewServer constructs a new Server using the specified config.
func NewServer(cfg Config) (*Server, error) {
	present.PlayEnabled = cfg.PlayEnabled
	s := &Server{cfg: cfg}

	parse := func(fs vfs.FileSystem, t *template.Template, filenames ...string) (*template.Template, error) {
		if t == nil {
			t = template.New(filenames[0]).Funcs(template.FuncMap{"playable": s.playable})
		} else {
			t = t.Funcs(template.FuncMap{"playable": s.playable})
		}
		for _, name := range filenames {
			data, err := vfs.ReadFile(fs, filepath.ToSlash(filepath.Join(cfg.TemplatePath, name)))
			if err != nil {
				return nil, err
			}
			if _, err := t.Parse(string(data)); err != nil {
				return nil, err
			}
		}
		return t, nil
	}

	// Parse templates.
	var err error
	s.template.action, err = parse(s.cfg.RootFS, nil, "action.tmpl", "action.tmpl")
	if err != nil {
		return nil, err
	}
	s.template.article, err = parse(s.cfg.RootFS, nil, "article.tmpl", "article.tmpl")
	if err != nil {
		return nil, err
	}
	s.template.dir, err = parse(s.cfg.RootFS, nil, "dir.tmpl", "dir.tmpl")
	if err != nil {
		return nil, err
	}
	s.template.slides, err = parse(s.cfg.RootFS, nil, "slides.tmpl", "slides.tmpl")
	if err != nil {
		return nil, err
	}
	s.template.doc, err = parse(s.cfg.RootFS, present.Template(), "doc.tmpl")
	if err != nil {
		return nil, err
	}

	s.contentTemplate = make(map[string]*template.Template)
	s.contentTemplate[".article"] = s.template.article
	s.contentTemplate[".slide"] = s.template.slides

	// Set up content file server.
	s.content = http.StripPrefix(s.cfg.BasePath, http.FileServer(
		httpfs.New(s.getNameSpace(s.cfg.RootFS, s.cfg.ContentPath)),
	))

	return s, nil
}

func (s *Server) getNameSpace(fs vfs.FileSystem, ns string) vfs.NameSpace {
	newns := make(vfs.NameSpace)
	if ns != "" {
		newns.Bind("/", fs, ns, vfs.BindReplace)
	} else {
		newns.Bind("/", fs, "/", vfs.BindReplace)
	}
	return newns
}

// ServeHTTP serves the front, index, and article pages
// as well as the ATOM and JSON feeds.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, s.cfg.BasePath)
	if s.isDoc(name) {
		err := s.renderDoc(w, name)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
		}
		return
	}
	if isDir, err := s.dirList(w, name); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	} else if isDir {
		return
	}

	s.content.ServeHTTP(w, r)
}

func (s *Server) playable(c present.Code) bool {
	return present.PlayEnabled && c.Play
}

func (s *Server) isDoc(path string) bool {
	_, ok := s.contentTemplate[filepath.Ext(path)]
	return ok
}

// renderDoc reads the present file, gets its template representation,
// and executes the template, sending output to w.
func (s *Server) renderDoc(w io.Writer, docFile string) error {
	// Read the input and build the doc structure.
	doc, err := s.parse(docFile, 0)
	if err != nil {
		return err
	}

	// Find which template should be executed.
	tmpl := s.contentTemplate[filepath.Ext(docFile)]

	// Execute the template.
	return doc.Render(w, tmpl)
}

func (s *Server) parse(name string, mode present.ParseMode) (*present.Doc, error) {
	f, err := s.cfg.RootFS.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return present.Parse(f, name, 0)
}

// dirList scans the given path and writes a directory listing to w.
// It parses the first part of each .slide file it encounters to display the
// presentation title in the listing.
// If the given path is not a directory, it returns (isDir == false, err == nil)
// and writes nothing to w.
func (s *Server) dirList(w io.Writer, name string) (isDir bool, err error) {
	f, err := s.cfg.RootFS.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	fi, err := s.cfg.RootFS.Stat(name)
	if err != nil {
		return false, err
	}
	if isDir = fi.IsDir(); !isDir {
		return false, nil
	}
	fis, err := s.cfg.RootFS.ReadDir(name)
	if err != nil {
		return false, err
	}
	d := &dirListData{Path: name}
	for _, fi := range fis {
		// skip the golang.org directory
		if name == "." && fi.Name() == "golang.org" {
			continue
		}
		e := dirEntry{
			Name: fi.Name(),
			Path: filepath.ToSlash(filepath.Join(name, fi.Name())),
		}
		if fi.IsDir() && s.showDir(e.Name) {
			d.Dirs = append(d.Dirs, e)
			continue
		}
		if s.isDoc(e.Name) {
			if p, err := s.parse(e.Path, present.TitlesOnly); err != nil {
				log.Println(err)
			} else {
				e.Title = p.Title
			}
			switch filepath.Ext(e.Path) {
			case ".article":
				d.Articles = append(d.Articles, e)
			case ".slide":
				d.Slides = append(d.Slides, e)
			}
		} else if s.showFile(e.Name) {
			d.Other = append(d.Other, e)
		}
	}
	if d.Path == "." {
		d.Path = ""
	}
	sort.Sort(d.Dirs)
	sort.Sort(d.Slides)
	sort.Sort(d.Articles)
	sort.Sort(d.Other)
	return true, s.template.dir.Execute(w, d)
}

// showFile reports whether the given file should be displayed in the list.
func (s *Server) showFile(n string) bool {
	switch filepath.Ext(n) {
	case ".pdf":
	case ".html":
	case ".go":
	default:
		return s.isDoc(n)
	}
	return true
}

// showDir reports whether the given directory should be displayed in the list.
func (s *Server) showDir(n string) bool {
	if len(n) > 0 && (n[0] == '.' || n[0] == '_') || n == "present" {
		return false
	}
	return true
}

type dirListData struct {
	Path                          string
	Dirs, Slides, Articles, Other dirEntrySlice
}

type dirEntry struct {
	Name, Path, Title string
}

type dirEntrySlice []dirEntry

func (s dirEntrySlice) Len() int           { return len(s) }
func (s dirEntrySlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s dirEntrySlice) Less(i, j int) bool { return s[i].Name < s[j].Name }
