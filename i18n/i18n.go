// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package i18n is for godoc Internationalization and Localization.

*/
package i18n

import (
	"fmt"
	"go/doc"

	"github.com/chai2010/golangdoc/godoc/static"
)

const (
	pkgId         = "__pkg__"
	pkgNameId     = "__name__"
	pkgSynopsisId = "__synopsis__"
	pkgDocId      = "__doc__"
)

var (
	staticFilesTable = make(map[string]map[string]string) // map[lang]...
	packageTable     = make(map[string]*doc.Package)      // map[mapKey(...)]...
	docTable         = make(map[string]string)            // map[mapKey(...)]...
)

func mapKey(lang, importPath, id string) string {
	return fmt.Sprintf("i18n:%s:%s", lang, importPath, id)
}

func MethodId(typeName, methodName string) string {
	return typeName + "." + methodName
}

func RegisterStaticFiles(lang string, staticFiles map[string]string) {
	staticFilesTable[lang] = staticFiles
}

func RegisterPackage(lang string, pkg *doc.Package) {
	packageTable[mapKey(lang, pkg.ImportPath, pkgId)] = pkg
	initDocTable(lang, pkg)
}

func initDocTable(lang string, pkg *doc.Package) {
	docTable[mapKey(lang, pkg.ImportPath, pkgNameId)] = pkg.Name
	docTable[mapKey(lang, pkg.ImportPath, pkgSynopsisId)] = doc.Synopsis(pkg.Doc)
	docTable[mapKey(lang, pkg.ImportPath, pkgDocId)] = pkg.Doc

	for _, v := range pkg.Consts {
		for _, id := range v.Names {
			docTable[mapKey(lang, pkg.ImportPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Types {
		docTable[mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc

		for _, x := range v.Consts {
			for _, id := range x.Names {
				docTable[mapKey(lang, pkg.ImportPath, id)] = x.Doc
			}
		}
		for _, x := range v.Vars {
			for _, id := range x.Names {
				docTable[mapKey(lang, pkg.ImportPath, id)] = x.Doc
			}
		}
		for _, x := range v.Funcs {
			docTable[mapKey(lang, pkg.ImportPath, x.Name)] = x.Doc
		}
		for _, x := range v.Methods {
			docTable[mapKey(lang, pkg.ImportPath, MethodId(v.Name, x.Name))] = x.Doc
		}
	}
	for _, v := range pkg.Vars {
		for _, id := range v.Names {
			docTable[mapKey(lang, pkg.ImportPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Funcs {
		docTable[mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc
	}
}

func StaticFiles(lang string) map[string]string {
	if f, ok := staticFilesTable[lang]; ok && f != nil {
		return f
	}
	return static.Files
}

func Synopsis(lang, importPath string, synopsis ...string) string {
	if s, ok := docTable[mapKey(lang, importPath, pkgSynopsisId)]; ok && s != "" {
		return s
	}
	if len(synopsis) > 0 {
		return synopsis[0]
	}
	return ""
}

func Doc(lang, importPath string, id ...string) (doc string) {
	if len(id) > 0 {
		doc, _ = docTable[mapKey(lang, importPath, id[0])]
	} else {
		doc, _ = docTable[mapKey(lang, importPath, pkgDocId)]
	}
	return
}

func Package(lang, importPath string, pkg ...*doc.Package) *doc.Package {
	if len(pkg) > 0 && pkg[0] != nil {
		return trPackage(lang, importPath, pkg[0])
	}
	if p, ok := packageTable[mapKey(lang, importPath, pkgId)]; ok && p != nil {
		return p
	}
	return nil
}

func trPackage(lang, importPath string, pkg *doc.Package) *doc.Package {
	key := mapKey(lang, pkg.ImportPath, pkgId)
	localPkg, _ := packageTable[key]
	if localPkg == nil {
		return pkg
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
		if s, _ := docTable[key]; s != "" {
			pkg.Consts[i].Doc = s
		}
	}
	for i := 0; i < len(pkg.Types); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Name)
		if s, _ := docTable[key]; s != "" {
			pkg.Types[i].Doc = s
		}

		for j := 0; j < len(pkg.Types[i].Consts); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Consts[j].Names[0])
			if s, _ := docTable[key]; s != "" {
				pkg.Types[i].Consts[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Vars); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Vars[j].Names[0])
			if s, _ := docTable[key]; s != "" {
				pkg.Types[i].Vars[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Funcs); j++ {
			key := mapKey(lang, pkg.ImportPath, pkg.Types[i].Funcs[j].Name)
			if s, _ := docTable[key]; s != "" {
				pkg.Types[i].Funcs[j].Doc = s
			}
		}
		for j := 0; j < len(pkg.Types[i].Methods); j++ {
			id := MethodId(pkg.Types[i].Name, pkg.Types[i].Methods[j].Name)
			key := mapKey(lang, pkg.ImportPath, id)
			if s, _ := docTable[key]; s != "" {
				pkg.Types[i].Methods[j].Doc = s
			}
		}
	}
	for i := 0; i < len(pkg.Vars); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Vars[i].Names[0])
		if s, _ := docTable[key]; s != "" {
			pkg.Vars[i].Doc = s
		}
	}
	for i := 0; i < len(pkg.Funcs); i++ {
		key := mapKey(lang, pkg.ImportPath, pkg.Funcs[i].Name)
		if s, _ := docTable[key]; s != "" {
			pkg.Funcs[i].Doc = s
		}
	}
	return pkg
}
