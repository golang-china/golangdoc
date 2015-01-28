// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"fmt"
	"go/doc"

	"golang.org/x/tools/godoc/vfs"
)

const (
	__pkg__  = "__pkg__"
	__name__ = "__name__"
	__doc__  = "__doc__"
)

// Translater interface.
type Translater interface {
	Static(lang string) vfs.FileSystem
	Document(lang string) vfs.FileSystem
	Package(lang, importPath string, pkg ...*doc.Package) *doc.Package
	Blog(lang string) vfs.FileSystem
}

var (
	staticFSTable    = make(map[string]vfs.FileSystem) // map[lang]...
	docFSTable       = make(map[string]vfs.FileSystem) // map[lang]...
	blogFSTable      = make(map[string]vfs.FileSystem) // map[lang]...
	pkgDocTable      = make(map[string]*doc.Package)   // map[mapKey(...)]...
	pkgDocIndexTable = make(map[string]string)         // map[mapKey(...)]...
	trList           = make([]Translater, 0)
)

func mapKey(lang, importPath, id string) string {
	return fmt.Sprintf("%s.%s@%s", importPath, id, lang)
}

func methodId(typeName, methodName string) string {
	return typeName + "." + methodName
}

// RegisterStaticFS Register StaticFS.
func RegisterStaticFS(lang string, staticFiles vfs.FileSystem) {
	staticFSTable[lang] = staticFiles
}

// RegisterDocumentFS Register DocumentFS.
func RegisterDocumentFS(lang string, docFiles vfs.FileSystem) {
	docFSTable[lang] = docFiles
}

// RegisterBlogFS Register BlogFS.
func RegisterBlogFS(lang string, blogFiles vfs.FileSystem) {
	blogFSTable[lang] = blogFiles
}

// RegisterPackage Register Package.
func RegisterPackage(lang string, pkg *doc.Package) {
	pkgDocTable[mapKey(lang, pkg.ImportPath, __pkg__)] = pkg
	initDocTable(lang, pkg)
}

// RegisterTranslater Register Translater.
func RegisterTranslater(tr Translater) {
	trList = append(trList, tr)
}

// RootFS return root filesystem.
func RootFS() vfs.FileSystem {
	return defaultRootFS
}

// StaticFS return Static filesystem.
func StaticFS(lang string) vfs.FileSystem {
	if lang == "" {
		return defaultStaticFS
	}
	if fs, _ := staticFSTable[lang]; fs != nil {
		return fs
	}
	for _, tr := range trList {
		if fs := tr.Static(lang); fs != nil {
			return fs
		}
	}
	if fs := defaultTranslater.Static(lang); fs != nil {
		return fs
	}
	return defaultStaticFS
}

// DocumentFS return Document filesystem.
func DocumentFS(lang string) vfs.FileSystem {
	if lang == "" {
		return defaultDocFS
	}
	if fs, _ := docFSTable[lang]; fs != nil {
		return fs
	}
	for _, tr := range trList {
		if fs := tr.Document(lang); fs != nil {
			return fs
		}
	}
	if fs := defaultTranslater.Document(lang); fs != nil {
		return fs
	}
	return defaultDocFS
}

// Package translate Package doc.
func Package(lang, importPath string, pkg ...*doc.Package) *doc.Package {
	if lang == "" {
		if len(pkg) > 0 {
			return pkg[0]
		} else {
			return nil
		}
	}
	if len(pkg) > 0 && pkg[0] != nil {
		if p := trPackage(lang, pkg[0].ImportPath, pkg[0]); p != nil {
			return p
		}
	} else {
		if p, _ := pkgDocTable[mapKey(lang, importPath, __pkg__)]; p != nil {
			return p
		}
	}
	for _, tr := range trList {
		if p := tr.Package(lang, importPath, pkg...); p != nil {
			return p
		}
	}
	if fs := defaultTranslater.Package(lang, importPath, pkg...); fs != nil {
		return fs
	}
	if len(pkg) > 0 {
		return pkg[0]
	}
	return nil
}

// BlogFS return Blog filesystem.
func BlogFS(lang string) vfs.FileSystem {
	if lang == "" {
		return defaultBlogFS
	}
	if fs, _ := blogFSTable[lang]; fs != nil {
		return fs
	}
	for _, tr := range trList {
		if fs := tr.Blog(lang); fs != nil {
			return fs
		}
	}
	if fs := defaultTranslater.Blog(lang); fs != nil {
		return fs
	}
	return defaultBlogFS
}

func initDocTable(lang string, pkg *doc.Package) {
	pkgDocIndexTable[mapKey(lang, pkg.ImportPath, __name__)] = pkg.Name
	pkgDocIndexTable[mapKey(lang, pkg.ImportPath, __doc__)] = pkg.Doc

	for _, v := range pkg.Consts {
		for _, id := range v.Names {
			pkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Types {
		pkgDocIndexTable[mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc

		for _, x := range v.Consts {
			for _, id := range x.Names {
				pkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = x.Doc
			}
		}
		for _, x := range v.Vars {
			for _, id := range x.Names {
				pkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = x.Doc
			}
		}
		for _, x := range v.Funcs {
			pkgDocIndexTable[mapKey(lang, pkg.ImportPath, x.Name)] = x.Doc
		}
		for _, x := range v.Methods {
			pkgDocIndexTable[mapKey(lang, pkg.ImportPath, methodId(v.Name, x.Name))] = x.Doc
		}
	}
	for _, v := range pkg.Vars {
		for _, id := range v.Names {
			pkgDocIndexTable[mapKey(lang, pkg.ImportPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Funcs {
		pkgDocIndexTable[mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc
	}
}

func trPackage(lang, importPath string, pkg *doc.Package) *doc.Package {
	key := mapKey(lang, pkg.ImportPath, __pkg__)
	localPkg, _ := pkgDocTable[key]
	if localPkg == nil {
		return nil
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
		if s, _ := pkgDocIndexTable[key]; s != "" {
			pkg.Consts[i].Doc = s
		}
	}
	for i := 0; i < len(pkg.Types); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Name)
		if s, _ := pkgDocIndexTable[key]; s != "" {
			pkg.Types[i].Doc = s
		}

		for j := 0; j < len(pkg.Types[i].Consts); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Consts[j].Names[0])
			if s, _ := pkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Consts[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Vars); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Vars[j].Names[0])
			if s, _ := pkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Vars[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Funcs); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Funcs[j].Name)
			if s, _ := pkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Funcs[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Methods); j++ {
			id := methodId(pkg.Types[i].Name, pkg.Types[i].Methods[j].Name)
			key := mapKey(lang, pkg.ImportPath, id)
			if s, _ := pkgDocIndexTable[key]; s != "" {
				pkg.Types[i].Methods[j].Doc = s
			}
		}
	}
	for i := 0; i < len(pkg.Vars); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Vars[i].Names[0])
		if s, _ := pkgDocIndexTable[key]; s != "" {
			pkg.Vars[i].Doc = s
		}
	}
	for i := 0; i < len(pkg.Funcs); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Funcs[i].Name)
		if s, _ := pkgDocIndexTable[key]; s != "" {
			pkg.Funcs[i].Doc = s
		}
	}
	return pkg
}
