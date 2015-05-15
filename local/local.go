// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"go/build"
	"go/doc"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/tools/godoc/static"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

// Default is the translations dir.
const (
	DefaultDir = "translations"     // $(RootFS)/translations
	DefaultEnv = "GODOC_LOCAL_ROOT" // dir list
)

var (
	defaultGodocGoos   = getGodocGoos()
	defaultGodocGoarch = getGodocGoarch()

	currentGoroot       string // runtime.GOROOT()
	currentGopath       string // build.Default.GOPATH
	currentTranslations string // os.Getenv(DefaultEnv)
)

var (
	gorootFS           vfs.NameSpace
	gorootFSTable      map[string]vfs.NameSpace // map[lang]...
	gopkgDocTable      map[string]*doc.Package  // map[mapKey(...)]...
	gopkgDocIndexTable map[string]string        // map[mapKey(...)]...
)

var (
	nilfs = make(vfs.NameSpace)
)

func init() {
	buildRootFS(runtime.GOROOT(), build.Default.GOPATH, "")
}

// Init initialize the translations environment.
func Init(goroot, gopath, translations string) {
	buildRootFS(goroot, gopath, translations)
}

func buildRootFS(goroot, gopath, translations string) {
	if gorootFS != nil && goroot == currentGoroot && gopath == currentGopath && translations == currentTranslations {
		return
	}

	currentGoroot = goroot
	currentGopath = gopath
	currentTranslations = translations

	gorootFSTable = make(map[string]vfs.NameSpace)
	gopkgDocTable = make(map[string]*doc.Package)
	gopkgDocIndexTable = make(map[string]string)

	var rootfs vfs.NameSpace
	switch {
	case strings.HasSuffix(goroot, ".zip"):
		rootfs = openZipFS(goroot)
	default:
		rootfs = getNameSpace(vfs.OS(runtime.GOROOT()), "/")
	}

	// gopath
	if gopath != "" {
		for _, p := range filepath.SplitList(gopath) {
			rootfs.Bind("/src", vfs.OS(p), "/src", vfs.BindAfter)
		}
	}

	// translations
	switch {
	case translations != "":
		if strings.HasSuffix(translations, ".zip") {
			rootfs.Bind("/translations", openZipFS(translations), "/", vfs.BindReplace)
		} else {
			rootfs.Bind("/translations", vfs.OS(translations), "/", vfs.BindReplace)
		}
	case os.Getenv(DefaultEnv) != "":
		for _, p := range filepath.SplitList(os.Getenv(DefaultEnv)) {
			fi, err := os.Lstat(p)
			if err != nil {
				log.Fatalf("local: os.Lstat(%q) failed: %s\n", p, err)
			}
			if strings.HasSuffix(fi.Name(), ".zip") {
				rootfs.Bind("/translations", openZipFS(p), "/", vfs.BindAfter)
			} else {
				rootfs.Bind("/translations", vfs.OS(p), "/", vfs.BindAfter)
			}
		}
	default:
		// default is `$(RootFS)/translations`
	}
	if _, err := rootfs.Lstat("/translations/src"); err == nil {
		rootfs.Bind("/src", rootfs, "/translations/src", vfs.BindAfter)
	}

	// lib/godoc
	if _, err := rootfs.Lstat("/lib/godoc"); err != nil {
		rootfs.Bind("/lib/godoc", mapfs.New(static.Files), "/", vfs.BindAfter)
	}

	// blog
	if _, err := rootfs.Lstat("/blog"); err != nil {
		const blogPath = "/src/golang.org/x/blog"
		if _, err := rootfs.Lstat(blogPath); err != nil {
			rootfs.Bind("/blog/static", getNameSpace(rootfs, blogPath), "/static", vfs.BindReplace)
			rootfs.Bind("/blog/template", getNameSpace(rootfs, blogPath), "/template", vfs.BindReplace)
			rootfs.Bind("/blog/content", getNameSpace(rootfs, blogPath), "/content", vfs.BindReplace)
		}
	}

	// talks
	if _, err := rootfs.Lstat("/talks"); err != nil {
		const presentPath = "/src/golang.org/x/tools/cmd/present"
		if _, err := rootfs.Lstat(presentPath); err != nil {
			rootfs.Bind("/talks/static", getNameSpace(rootfs, presentPath), "/static", vfs.BindReplace)
			rootfs.Bind("/talks/template", getNameSpace(rootfs, presentPath), "/template", vfs.BindReplace)
		}
		const talksPath = "golang.org/x/talks"
		if _, err := rootfs.Lstat("/src/" + talksPath); err != nil {
			rootfs.Bind("/talks/content", getNameSpace(rootfs, talksPath), talksPath, vfs.BindReplace)
		}
	}

	// tour
	if _, err := rootfs.Lstat("/tour"); err != nil {
		const tourPath = "golang.org/x/tour"
		if _, err := rootfs.Lstat(tourPath); err != nil {
			rootfs.Bind("/tour/static", getNameSpace(rootfs, tourPath), tourPath+"/static", vfs.BindReplace)
			rootfs.Bind("/tour/template", getNameSpace(rootfs, tourPath), tourPath+"/template", vfs.BindReplace)
			rootfs.Bind("/tour/content", getNameSpace(rootfs, tourPath), tourPath+"/content", vfs.BindReplace)
		}
	}

	gorootFS = rootfs
	return
}

// RootFS return root filesystem.
func RootFS(lang string) vfs.NameSpace {
	if lang == "" {
		return gorootFS
	}
	if fs, _ := gorootFSTable[lang]; fs != nil {
		return fs
	}

	newfs := getNameSpace(gorootFS, "/")
	{
		// lib/godoc
		if _, err := gorootFS.Lstat("/translations/static/" + lang); err == nil {
			newfs.Bind("/lib/godoc", gorootFS, "/translations/static/"+lang, vfs.BindReplace)
		}

		// doc
		if _, err := gorootFS.Lstat("/translations/doc/" + lang); err == nil {
			newfs.Bind("/doc", gorootFS, "/translations/doc/"+lang, vfs.BindReplace)
		}

		// blog
		if _, err := gorootFS.Lstat("/translations/blog/" + lang); err == nil {
			newfs.Bind("/blog", gorootFS, "/translations/blog/"+lang, vfs.BindReplace)
		}

		// talks
		if _, err := gorootFS.Lstat("/translations/talks/" + lang); err == nil {
			newfs.Bind("/talks", gorootFS, "/translations/talks/"+lang, vfs.BindReplace)
		}

		// tour
		if _, err := gorootFS.Lstat("/translations/tour/" + lang); err == nil {
			newfs.Bind("/tour", gorootFS, "/translations/tour/"+lang, vfs.BindReplace)
		}
	}

	gorootFSTable[lang] = newfs
	return newfs
}

// Package translate Package doc.
func Package(lang, importPath string, in *doc.Package) *doc.Package {
	key := mapKey(lang, importPath, __pkg__)

	// build package doc
	if _, ok := gopkgDocTable[key]; !ok {
		pkg := parsePkgDocPackage(RootFS(lang), lang, importPath)
		gopkgDocTable[key] = pkg
		if pkg == nil {
			return in
		}

		gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, __name__)] = pkg.Name
		gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, __doc__)] = pkg.Doc

		for _, v := range pkg.Consts {
			for _, id := range v.Names {
				gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = v.Doc
			}
		}
		for _, v := range pkg.Types {
			gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc

			for _, x := range v.Consts {
				for _, id := range x.Names {
					gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = x.Doc
				}
			}
			for _, x := range v.Vars {
				for _, id := range x.Names {
					gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = x.Doc
				}
			}
			for _, x := range v.Funcs {
				gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, x.Name)] = x.Doc
			}
			for _, x := range v.Methods {
				gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, methodId(v.Name, x.Name))] = x.Doc
			}
		}
		for _, v := range pkg.Vars {
			for _, id := range v.Names {
				gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = v.Doc
			}
		}
		for _, v := range pkg.Funcs {
			gopkgDocIndexTable[mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc
		}
	}

	return trPackage(lang, importPath, in)
}

func trPackage(lang, importPath string, pkg *doc.Package) *doc.Package {
	key := mapKey(lang, importPath, __pkg__)
	localPkg, _ := gopkgDocTable[key]
	if localPkg == nil {
		return pkg
	}
	if pkg == nil {
		return localPkg
	}

	pkg.Name = localPkg.Name
	pkg.Doc = localPkg.Doc

	for k, _ := range pkg.Notes {
		if notes, _ := localPkg.Notes[k]; notes != nil {
			pkg.Notes[k] = notes
		}
	}

	for i := 0; i < len(pkg.Consts); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Consts[i].Names[0])
		if s, _ := gopkgDocIndexTable[key]; s != "" {
			pkg.Consts[i].Doc = s
		}
	}
	for i := 0; i < len(pkg.Types); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Name)
		if s, _ := gopkgDocIndexTable[key]; s != "" {
			pkg.Types[i].Doc = s
		}

		for j := 0; j < len(pkg.Types[i].Consts); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Consts[j].Names[0])
			if s, _ := gopkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Consts[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Vars); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Vars[j].Names[0])
			if s, _ := gopkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Vars[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Funcs); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Funcs[j].Name)
			if s, _ := gopkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Funcs[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Methods); j++ {
			id := methodId(pkg.Types[i].Name, pkg.Types[i].Methods[j].Name)
			key := mapKey(lang, pkg.ImportPath, id)
			if s, _ := gopkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Methods[j].Doc = s
			}
		}
	}
	for i := 0; i < len(pkg.Vars); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Vars[i].Names[0])
		if s, _ := gopkgDocIndexTable[key]; s != "" {
			pkg.Vars[i].Doc = s
		}
	}
	for i := 0; i < len(pkg.Funcs); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Funcs[i].Name)
		if s, _ := gopkgDocIndexTable[key]; s != "" {
			pkg.Funcs[i].Doc = s
		}
	}
	return pkg
}
