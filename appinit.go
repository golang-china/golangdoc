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

	"golang.org/x/tools/godoc/vfs"

	"github.com/golang-china/golangdoc/godoc"
	"github.com/golang-china/golangdoc/local"
)

func init() {
	playEnabled = true

	log.Println("initializing godoc ...")
	log.Printf(".zip file   = %s", zipFilename)
	log.Printf(".zip GOROOT = %s", zipGoroot)
	log.Printf("index files = %s", indexFilenames)

	// Determine file system to use.
	local.Init(zipGoroot, zipFilename, "", "")
	fs.Bind("/", local.RootFS(), "/", vfs.BindReplace)
	fs.Bind("/lib/godoc", local.StaticFS(*lang), "/", vfs.BindReplace)
	fs.Bind("/doc", local.DocumentFS(*lang), "/", vfs.BindReplace)

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = false
	corpus.MaxResults = 10000 // matches flag default in main.go
	corpus.IndexEnabled = true
	corpus.IndexFiles = indexFilenames

	// translate hook
	corpus.SummarizePackage = func(importPath string) (summary string, showList, ok bool) {
		if pkg := local.Package(*lang, importPath); pkg != nil {
			summary = doc.Synopsis(pkg.Doc)
		}
		ok = (summary != "")
		return
	}
	corpus.TranslateDocPackage = func(pkg *doc.Package) *doc.Package {
		return local.Package(*lang, pkg.ImportPath, pkg)
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
