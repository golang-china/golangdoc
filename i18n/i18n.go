// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package i18n is for godoc Internationalization and Localization.

*/
package i18n

import (
	"go/doc"

	"golang.org/x/tools/godoc/static"
)

type local struct {
	lang             string
	localDocPackage  func(p *doc.Package) *doc.Package
	localStaticFiles func() map[string]string
}

var locals []local

func Register(lang string,
	localDocPackage func(p *doc.Package) *doc.Package,
	localStaticFiles func() map[string]string,
) {
	locals = append(locals, local{
		lang:             lang,
		localDocPackage:  localDocPackage,
		localStaticFiles: localStaticFiles,
	})
}

func sniff(lang string) local {
	for _, f := range locals {
		if f.lang == lang {
			return f
		}
	}
	return local{}
}

func DocPackage(p *doc.Package, lang ...string) *doc.Package {
	if len(lang) > 0 {
		if f := sniff(lang[0]); f.localDocPackage != nil {
			return f.localDocPackage(p)
		}
	}
	return p
}

func StaticFiles(lang ...string) map[string]string {
	if len(lang) > 0 {
		if f := sniff(lang[0]); f.localStaticFiles != nil {
			return f.localStaticFiles()
		}
	}
	return static.Files
}
