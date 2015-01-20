// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run generate.go

package pkgdoc

import (
	"go/doc"
)

func OriginDocPackage(lang, importPath string) *doc.Package {
	pkgList, _ := OriginDocPackageTable[lang]
	for i := 0; i < len(pkgList); i++ {
		if importPath == pkgList[i].ImportPath {
			return pkgList[i]
		}
	}
	return nil
}

func TranslateDocPackage(lang, importPath string) *doc.Package {
	pkgList, _ := TranslateDocPackageTable[lang]
	for i := 0; i < len(pkgList); i++ {
		if importPath == pkgList[i].ImportPath {
			return pkgList[i]
		}
	}
	return nil
}
