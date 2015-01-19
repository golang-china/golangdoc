// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run generate.go

package pkgdoc

import (
	"fmt"
	"go/doc"
)

const (
	PackageId    = "" // Same as synopsis
	SynopsisId   = "__synopsis__"
	PackageDocId = "__doc__"
)

var (
	locakPkg = map[string]*doc.Package{}
	locakDoc = map[string]string{}
)

func mapKey(lang, name, importPath, id string) string {
	return fmt.Sprintf("pkgdoc[%s]:%s:%s.%s", lang, name, importPath, id)
}

func MethodId(typeName, methodName string) string {
	return typeName + "." + methodName
}

func GetSynopsis(lang, name, importPath string) (doc string) {
	return GetDoc(lang, name, importPath, SynopsisId)
}

func GetDoc(lang, name, importPath string, id ...string) (doc string) {
	if len(id) > 0 {
		doc, _ = locakDoc[mapKey(lang, name, importPath, id[0])]
	} else {
		doc, _ = locakDoc[mapKey(lang, name, importPath, "")]
	}
	return
}

func GetPackage(lang, name, importPath string, pkg ...*doc.Package) *doc.Package {
	pkg1, _ := locakPkg[mapKey(lang, name, importPath, "")]
	return pkg1
}

func Register(lang string, pkg *doc.Package) {
	name, importPath := pkg.Name, pkg.ImportPath

	locakPkg[mapKey(lang, name, importPath, "")] = pkg

	locakDoc[mapKey(lang, name, importPath, PackageId)] = doc.Synopsis(pkg.Doc)
	locakDoc[mapKey(lang, name, importPath, SynopsisId)] = doc.Synopsis(pkg.Doc)
	locakDoc[mapKey(lang, name, importPath, PackageDocId)] = pkg.Doc

	for _, v := range pkg.Consts {
		for _, id := range v.Names {
			locakDoc[mapKey(lang, name, importPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Types {
		locakDoc[mapKey(lang, name, importPath, v.Name)] = v.Doc

		for _, x := range v.Consts {
			for _, id := range x.Names {
				locakDoc[mapKey(lang, name, importPath, id)] = x.Doc
			}
		}
		for _, x := range v.Vars {
			for _, id := range x.Names {
				locakDoc[mapKey(lang, name, importPath, id)] = x.Doc
			}
		}
		for _, x := range v.Funcs {
			locakDoc[mapKey(lang, name, importPath, x.Name)] = x.Doc
		}
		for _, x := range v.Methods {
			locakDoc[mapKey(lang, name, importPath, MethodId(v.Name, x.Name))] = x.Doc
		}
	}
	for _, v := range pkg.Vars {
		for _, id := range v.Names {
			locakDoc[mapKey(lang, name, importPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Funcs {
		locakDoc[mapKey(lang, name, importPath, v.Name)] = v.Doc
	}
}
