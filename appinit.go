// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

package main

// This file replaces main.go when running godoc under app-engine.
// See README.godoc-app for details.

import (
	"go/doc"
	"log"
	"regexp"

	"github.com/golang-china/golangdoc/godoc"
	"github.com/golang-china/golangdoc/local"
)

func init() {
	playEnabled = true

	log.Println("initializing godoc ...")
	log.Printf(".zip GOROOT = %s", flagZipGoroot)
	log.Printf("index files = %s", flagIndexFilenames)

	// Determine file system to use.
	local.Init(flagZipGoroot, "", "")
	fs = local.RootFS(*flagLang)

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = false
	corpus.MaxResults = 10000 // matches flag default in main.go
	corpus.IndexEnabled = true
	corpus.IndexFiles = flagIndexFilenames

	// translate hook
	corpus.SummarizePackage = func(importPath string, langs ...string) (summary string, showList, ok bool) {
		lang := *flagLang
		if len(langs) > 0 && langs[0] != "" {
			lang = langs[0]
		}
		if lang == "en" || lang == "raw" || lang == "EN" {
			lang = ""
		}
		if pkg := local.Package(lang, importPath, nil); pkg != nil {
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

	if err := corpus.Init(); err != nil {
		log.Fatal(err)
	}
	if corpus.IndexEnabled && corpus.IndexFiles != "" {
		go corpus.RunIndexer()
	}

	pres = godoc.NewPresentation(corpus)
	pres.TabWidth = 8
	pres.ShowPlayground = true
	pres.ShowExamples = true
	pres.DeclLinks = true
	pres.NotesRx = regexp.MustCompile("BUG")

	readTemplates(pres, true)
	registerHandlers(pres)

	log.Println("godoc initialization complete")
}
