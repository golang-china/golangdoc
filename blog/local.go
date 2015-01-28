// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"log"
	"net/http"
	"path"
	"runtime"

	"golang.org/x/tools/godoc/static"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"

	"github.com/chai2010/golangdoc/blog"
)

const (
	hostname = "blog.golang.org"
)

var cfg = blog.Config{
	Hostname:     hostname,
	RootFS:       vfs.OS(path.Join(runtime.GOROOT(), `/translations/blog/zh_CN`)),
	ContentPath:  "content",
	TemplatePath: "template",
	BaseURL:      "//" + hostname,
	GodocURL:     "//golang.org",
	HomeArticles: 5,  // articles to display on the home page
	FeedArticles: 10, // articles to include in Atom and JSON feeds
	PlayEnabled:  true,
	FeedTitle:    "The Go Programming Language Blog",
}

func main() {
	server, err := blog.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", server)
	http.Handle("/lib/godoc/", http.StripPrefix("/lib/godoc/",
		http.FileServer(httpfs.New(mapfs.New(static.Files))),
	))

	log.Fatal(http.ListenAndServe(":3999", nil))
}
