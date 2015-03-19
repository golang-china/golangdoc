// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// godoc: Go Documentation Server

// Web server tree:
//
//	http://godoc/		main landing page
//	http://godoc/doc/	serve from $GOROOT/doc - spec, mem, etc.
//	http://godoc/src/	serve files from $GOROOT/src; .go gets pretty-printed
//	http://godoc/cmd/	serve documentation about commands
//	http://godoc/pkg/	serve documentation about packages
//				(idea is if you say import "compress/zlib", you go to
//				http://godoc/pkg/compress/zlib)
//
// Command-line interface:
//
//	godoc packagepath [name ...]
//
//	godoc compress/zlib
//		- prints doc for package compress/zlib
//	godoc crypto/block Cipher NewCMAC
//		- prints doc for Cipher and NewCMAC in package crypto/block

// +build !appengine

package main

import (
	_ "expvar" // to serve /debug/vars
	"flag"
	"fmt"
	"go/build"
	"go/doc"
	"log"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof" // to serve /debug/pprof/*
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/tools/godoc/analysis"
	"golang.org/x/tools/godoc/vfs"

	"github.com/golang-china/golangdoc/godoc"
	"github.com/golang-china/golangdoc/local"
)

const (
	defaultAddr = ":6060" // default webserver address
	toolsPath   = "golang.org/x/tools/cmd/"
)

var (
	// file system to serve
	// (with e.g.: zip -r go.zip $GOROOT -i \*.go -i \*.html -i \*.css -i \*.js -i \*.txt -i \*.c -i \*.h -i \*.s -i \*.png -i \*.jpg -i \*.sh -i favicon.ico)
	flagZipfile = flag.String("zip", "", "zip file providing the file system to serve; disabled if empty")

	// file-based index
	flagWriteIndex = flag.Bool("write_index", false, "write index to a file; the file name must be specified with -index_files")

	flagAnalysisFlag = flag.String("analysis", "", `comma-separated list of analyses to perform (supported: type, pointer). See http://golang.org/lib/godoc/analysis/help.html`)

	// network
	flagHttpAddr   = flag.String("http", "", "HTTP service address (e.g., '"+defaultAddr+"')")
	flagServerAddr = flag.String("server", "", "webserver address for command line searches")

	// layout control
	flagHtml    = flag.Bool("html", false, "print HTML in command-line mode")
	flagSrcMode = flag.Bool("src", false, "print (exported) source in command-line mode")
	flagUrlFlag = flag.String("url", "", "print HTML for named URL")

	// command-line searches
	flagQuery = flag.Bool("q", false, "arguments are considered search queries")

	flagVerbose = flag.Bool("v", false, "verbose mode")

	// file system roots
	// TODO(gri) consider the invariant that goroot always end in '/'
	flagGoroot    = flag.String("goroot", runtime.GOROOT(), "Go root directory")
	flagLocalRoot = flag.String("godoc-local-root", "", "Godoc translates root, default is $(GOROOT)/translates")

	// layout control
	flagTabWidth       = flag.Int("tabwidth", 4, "tab width")
	flagShowTimestamps = flag.Bool("timestamps", false, "show timestamps with directory listings")
	flagTemplateDir    = flag.String("templates", "", "directory containing alternate template files")
	flagShowPlayground = flag.Bool("play", false, "enable playground in web interface")
	flagShowExamples   = flag.Bool("ex", false, "show examples in command line mode")
	flagDeclLinks      = flag.Bool("links", true, "link identifiers to their declarations")

	// search index
	flagIndexEnabled  = flag.Bool("index", false, "enable search index")
	flagIndexFiles    = flag.String("index_files", "", "glob pattern specifying index files; if not empty, the index is read from these files in sorted order")
	flagMaxResults    = flag.Int("maxresults", 10000, "maximum number of full text search results shown")
	flagIndexThrottle = flag.Float64("index_throttle", 0.75, "index throttle value; 0.0 = no time allocated, 1.0 = full throttle")

	// source code notes
	flagNotesRx = flag.String("notes", "BUG", "regular expression matching note markers to show")

	// local language
	flagLang = flag.String("lang", "", "local language")
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"usage: golangdoc package [name ...]\n"+
			"	golangdoc -http="+defaultAddr+"\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s\t%s", req.RemoteAddr, req.URL)
		h.ServeHTTP(w, req)
	})
}

func handleURLFlag() {
	// Try up to 10 fetches, following redirects.
	urlstr := *flagUrlFlag
	for i := 0; i < 10; i++ {
		// Prepare request.
		u, err := url.Parse(urlstr)
		if err != nil {
			log.Fatal(err)
		}
		req := &http.Request{
			URL: u,
		}

		// Invoke default HTTP handler to serve request
		// to our buffering httpWriter.
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, req)

		// Return data, error, or follow redirect.
		switch w.Code {
		case 200: // ok
			os.Stdout.Write(w.Body.Bytes())
			return
		case 301, 302, 303, 307: // redirect
			redirect := w.HeaderMap.Get("Location")
			if redirect == "" {
				log.Fatalf("HTTP %d without Location header", w.Code)
			}
			urlstr = redirect
		default:
			log.Fatalf("HTTP error %d", w.Code)
		}
	}
	log.Fatalf("too many redirects")
}

func runGodoc() {
	if *flagLocalRoot == "" {
		if s := os.Getenv("GODOC_LOCAL_ROOT"); s != "" {
			*flagLocalRoot = s
		}
	}

	// Determine file system to use.
	local.Init(*flagGoroot, *flagLocalRoot, *flagZipfile, *flagTemplateDir, build.Default.GOPATH)
	fs.Bind("/", local.RootFS(), "/", vfs.BindReplace)
	fs.Bind("/lib/godoc", local.StaticFS(*flagLang), "/", vfs.BindReplace)
	fs.Bind("/doc", local.DocumentFS(*flagLang), "/", vfs.BindReplace)

	httpMode := *flagHttpAddr != ""

	var typeAnalysis, pointerAnalysis bool
	if *flagAnalysisFlag != "" {
		for _, a := range strings.Split(*flagAnalysisFlag, ",") {
			switch a {
			case "type":
				typeAnalysis = true
			case "pointer":
				pointerAnalysis = true
			default:
				log.Fatalf("unknown analysis: %s", a)
			}
		}
	}

	corpus := godoc.NewCorpus(fs)

	// translate hook
	corpus.SummarizePackage = func(importPath string, langs ...string) (summary string, showList, ok bool) {
		lang := *flagLang
		if len(langs) > 0 && langs[0] != "" {
			lang = langs[0]
		}
		if lang == "en" || lang == "raw" || lang == "EN" {
			lang = ""
		}
		if pkg := local.Package(lang, importPath); pkg != nil {
			summary = doc.Synopsis(pkg.Doc)
		}
		ok = (summary != "")
		return
	}
	corpus.TranslateDocPackage = func(pkg *doc.Package, langs ...string) *doc.Package {
		lang := *flagLang
		if len(langs) > 0 && langs[0] != "" {
			lang = langs[0]
		}
		if lang == "en" || lang == "raw" || lang == "EN" {
			lang = ""
		}
		return local.Package(lang, pkg.ImportPath, pkg)
	}

	corpus.Verbose = *flagVerbose
	corpus.MaxResults = *flagMaxResults
	corpus.IndexEnabled = *flagIndexEnabled && httpMode
	if *flagMaxResults == 0 {
		corpus.IndexFullText = false
	}
	corpus.IndexFiles = *flagIndexFiles
	corpus.IndexThrottle = *flagIndexThrottle
	if *flagWriteIndex {
		corpus.IndexThrottle = 1.0
		corpus.IndexEnabled = true
	}
	if *flagWriteIndex || httpMode || *flagUrlFlag != "" {
		if err := corpus.Init(); err != nil {
			log.Fatal(err)
		}
	}

	pres = godoc.NewPresentation(corpus)
	pres.TabWidth = *flagTabWidth
	pres.ShowTimestamps = *flagShowTimestamps
	pres.ShowPlayground = *flagShowPlayground
	pres.ShowExamples = *flagShowExamples
	pres.DeclLinks = *flagDeclLinks
	pres.SrcMode = *flagSrcMode
	pres.HTMLMode = *flagHtml
	if *flagNotesRx != "" {
		pres.NotesRx = regexp.MustCompile(*flagNotesRx)
	}

	readTemplates(pres, httpMode || *flagUrlFlag != "")
	registerHandlers(pres)

	if *flagWriteIndex {
		// Write search index and exit.
		if *flagIndexFiles == "" {
			log.Fatal("no index file specified")
		}

		log.Println("initialize file systems")
		*flagVerbose = true // want to see what happens

		corpus.UpdateIndex()

		log.Println("writing index file", *flagIndexFiles)
		f, err := os.Create(*flagIndexFiles)
		if err != nil {
			log.Fatal(err)
		}
		index, _ := corpus.CurrentIndex()
		_, err = index.WriteTo(f)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("done")
		return
	}

	// Print content that would be served at the URL *urlFlag.
	if *flagUrlFlag != "" {
		handleURLFlag()
		return
	}

	if httpMode {
		// HTTP server mode.
		var handler http.Handler = http.DefaultServeMux
		if *flagVerbose {
			log.Printf("Go Documentation Server")
			log.Printf("version = %s", runtime.Version())
			log.Printf("address = %s", *flagHttpAddr)
			log.Printf("goroot = %s", *flagGoroot)
			log.Printf("tabwidth = %d", *flagTabWidth)
			switch {
			case !*flagIndexEnabled:
				log.Print("search index disabled")
			case *flagMaxResults > 0:
				log.Printf("full text index enabled (maxresults = %d)", *flagMaxResults)
			default:
				log.Print("identifier search index enabled")
			}
			handler = loggingHandler(handler)
		}

		// Initialize search index.
		if *flagIndexEnabled {
			go corpus.RunIndexer()
		}

		// Start type/pointer analysis.
		if typeAnalysis || pointerAnalysis {
			go analysis.Run(pointerAnalysis, &corpus.Analysis)
		}

		// Start http server.
		if err := http.ListenAndServe(*flagHttpAddr, handler); err != nil {
			log.Fatalf("ListenAndServe %s: %v", *flagHttpAddr, err)
		}

		return
	}

	if *flagQuery {
		handleRemoteSearch()
		return
	}

	if err := godoc.CommandLine(os.Stdout, fs, pres, flag.Args(), *flagLang); err != nil {
		log.Print(err)
	}
}
