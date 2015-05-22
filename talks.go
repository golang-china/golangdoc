// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/golang-china/golangdoc/internal/redirect"
	"github.com/golang-china/golangdoc/internal/talks"
	"github.com/golang-china/golangdoc/local"
)

const (
	talksRepo = "golang.org/x/talks"
	talksURL  = "http://talks.golang.org/"
	talksPath = "/talks/"
)

var (
	talksServer   http.Handler // set by talksInit
	talksInitOnce sync.Once
)

func init() {
	return

	// Initialize talks only when first accessed.
	http.HandleFunc(talksPath, func(w http.ResponseWriter, r *http.Request) {
		talksInitOnce.Do(talksInit)
		talksServer.ServeHTTP(w, r)
	})
}

func talksInit() {
	talksFS := getNameSpace(local.RootFS(*flagLang), "/talks")

	// If content is not available fall back to redirect.
	if fi, err := talksFS.Lstat("/"); err != nil || !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "Talks content not available locally. "+
			"To install, run \n\tgo get %v\n", talksRepo)
		talksServer = http.HandlerFunc(talksRedirectHandler)
		return
	}

	s, err := talks.NewServer(talks.Config{
		RootFS:       talksFS,
		BaseURL:      talksPath,
		BasePath:     strings.TrimSuffix(talksPath, "/"),
		ContentPath:  "content",
		TemplatePath: "template",
		HomeArticles: 5,
		PlayEnabled:  playEnabled,
	})
	if err != nil {
		log.Fatal(err)
	}
	talksServer = s
}

func talksRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == talksPath {
		http.Redirect(w, r, talksURL, http.StatusFound)
		return
	}
	talksPrefixHandler.ServeHTTP(w, r)
}

var talksPrefixHandler = redirect.PrefixHandler(talksPath, talksURL)
