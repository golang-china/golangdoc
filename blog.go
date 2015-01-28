// Copyright 2013 The Go Authors. All rights reserved.
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

	"golang.org/x/tools/godoc/redirect"

	"github.com/chai2010/golangdoc/blog"
	"github.com/chai2010/golangdoc/local"
)

const (
	blogRepo = "golang.org/x/blog"
	blogURL  = "http://blog.golang.org/"
	blogPath = "/blog/"
)

var (
	blogServer   http.Handler // set by blogInit
	blogInitOnce sync.Once
	playEnabled  bool
)

func init() {
	// Initialize blog only when first accessed.
	http.HandleFunc(blogPath, func(w http.ResponseWriter, r *http.Request) {
		blogInitOnce.Do(blogInit)
		blogServer.ServeHTTP(w, r)
	})
}

func blogInit() {
	blogFS := local.BlogFS(*lang)

	// If content is not available fall back to redirect.
	if fi, err := blogFS.Lstat("/"); err != nil || !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "Blog content not available locally. "+
			"To install, run \n\tgo get %v\n", blogRepo)
		blogServer = http.HandlerFunc(blogRedirectHandler)
		return
	}

	s, err := blog.NewServer(blog.Config{
		RootFS:       blogFS,
		BaseURL:      blogPath,
		BasePath:     strings.TrimSuffix(blogPath, "/"),
		ContentPath:  "content",
		TemplatePath: "template",
		HomeArticles: 5,
		PlayEnabled:  playEnabled,
	})
	if err != nil {
		log.Fatal(err)
	}
	blogServer = s
}

func blogRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == blogPath {
		http.Redirect(w, r, blogURL, http.StatusFound)
		return
	}
	blogPrefixHandler.ServeHTTP(w, r)
}

var blogPrefixHandler = redirect.PrefixHandler(blogPath, blogURL)
