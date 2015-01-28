// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"archive/zip"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/tools/godoc/static"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

// Default is the translations dir.
const (
	Default = "translations" // $(RootFS)/translations
)

var (
	defaultGodocGoos                 = getGodocGoos()
	defaultGodocGoarch               = getGodocGoarch()
	defaultRootFS      vfs.NameSpace = getNameSpace(vfs.OS(runtime.GOROOT()), "/")
	defaultStaticFS    vfs.NameSpace = getNameSpace(mapfs.New(static.Files), "/")
	defaultDocFS       vfs.NameSpace = getNameSpace(defaultRootFS, "/doc")
	defaultBlogFS      vfs.NameSpace = getNameSpace(defaultRootFS, "/blog")
	defaultLocalFS     vfs.NameSpace = getNameSpace(defaultRootFS, "/"+Default)
	defaultTranslater  Translater    = new(localTranslater)
)

func getGodocGoos() string {
	if v := strings.TrimSpace(os.Getenv("GOOS")); v != "" {
		return v
	}
	return runtime.GOOS
}

func getGodocGoarch() string {
	if v := strings.TrimSpace(os.Getenv("GOARCH")); v != "" {
		return v
	}
	return runtime.GOARCH
}

// Init initialize the translations environment.
func Init(goRoot, goZipFile, goTemplateDir, goPath string) {
	if goZipFile != "" {
		rc, err := zip.OpenReader(goZipFile)
		if err != nil {
			log.Fatalf("local: %s: %s\n", goZipFile, err)
		}
		defer rc.Close()

		defaultRootFS = getNameSpace(zipfs.New(rc, goZipFile), goRoot)
		defaultDocFS = getNameSpace(defaultRootFS, "/doc")
		defaultBlogFS = getNameSpace(defaultRootFS, "/blog")
		defaultLocalFS = getNameSpace(defaultRootFS, "/"+Default)
	} else {
		if goRoot != "" && goRoot != runtime.GOROOT() {
			defaultRootFS = getNameSpace(vfs.OS(goRoot), "/")
			defaultDocFS = getNameSpace(defaultRootFS, "/doc")
			defaultBlogFS = getNameSpace(defaultRootFS, "/blog")
			defaultLocalFS = getNameSpace(defaultRootFS, "/"+Default)
		}
	}

	if goTemplateDir != "" {
		defaultStaticFS = getNameSpace(vfs.OS(goTemplateDir), "/")
	}

	// Bind $GOPATH trees into Go root.
	for _, p := range filepath.SplitList(goPath) {
		defaultRootFS.Bind("/src", vfs.OS(p), "/src", vfs.BindAfter)
	}

	// Prefer content from go.blog repository if present.
	if _, err := defaultBlogFS.Lstat("/"); err != nil {
		const blogRepo = "golang.org/x/blog"
		if pkg, err := build.Import(blogRepo, "", build.FindOnly); err == nil {
			defaultBlogFS = getNameSpace(defaultRootFS, pkg.Dir)
		}
	}
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
