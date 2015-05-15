// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"flag"
	"log"
	"net/http"

	"golang.org/x/tools/godoc/static"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"

	"github.com/golang-china/golangdoc/blog"
	"github.com/golang-china/golangdoc/local"
)

const (
	hostname = "blog.golang.org"
)

var (
	flagHttpAddr = flag.String("http", ":3999", "HTTP service address")
	flagLang     = flag.String("lang", "zh_CN", "local language")
)

func main() {
	flag.Parse()

	rootfs := local.RootFS(*flagLang)
	blogfs := getNameSpace(rootfs, "/blog")

	var cfg = blog.Config{
		Hostname:     hostname,
		RootFS:       blogfs,
		ContentPath:  "content",
		TemplatePath: "template",
		BaseURL:      "//" + hostname,
		GodocURL:     "//golang.org",
		HomeArticles: 5,  // articles to display on the home page
		FeedArticles: 10, // articles to include in Atom and JSON feeds
		PlayEnabled:  true,
		FeedTitle:    "The Go Programming Language Blog",
	}

	server, err := blog.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", server)
	http.Handle("/lib/godoc/", http.StripPrefix("/lib/godoc/",
		http.FileServer(httpfs.New(mapfs.New(static.Files))),
	))

	log.Fatal(http.ListenAndServe(*flagHttpAddr, nil))
}

func getNameSpace(fs vfs.FileSystem, ns string) vfs.NameSpace {
	newns := make(vfs.NameSpace)
	if ns != "" {
		newns.Bind("/", fs, ns, vfs.BindReplace)
	} else {
		newns.Bind("/", fs, "/", vfs.BindReplace)
	}
	return newns
}
