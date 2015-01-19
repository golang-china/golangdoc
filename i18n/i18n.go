// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package i18n is for godoc Internationalization and Localization.

*/
package i18n

import (
	"go/doc"

	"github.com/chai2010/golangdoc/godoc/static"
)

var (
	localStaticFilesTable = map[string]map[string]string{}
	localDocPackageTable  = map[string]func(p *doc.Package) *doc.Package{}
)

func RegisterStaticFiles(lang string, localStaticFiles map[string]string) {
	localStaticFilesTable[lang] = localStaticFiles
}

func RegisterDocPackage(lang string, localDocPackage func(p *doc.Package) *doc.Package) {
	localDocPackageTable[lang] = localDocPackage
}

func StaticFiles(lang ...string) map[string]string {
	if len(lang) > 0 {
		if f, ok := localStaticFilesTable[lang[0]]; ok && f != nil {
			return f
		}
	}
	return static.Files
}

func DocPackage(p *doc.Package, lang ...string) *doc.Package {
	if len(lang) > 0 {
		if f, ok := localDocPackageTable[lang[0]]; ok && f != nil {
			if x := f(p); x != nil {
				return x
			}
		}
	}
	return p
}
