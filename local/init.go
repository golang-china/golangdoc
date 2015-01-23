// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"archive/zip"
	"log"
	"path/filepath"
	"runtime"

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
	defaultRootFS     vfs.NameSpace = getNameSpace(vfs.OS(runtime.GOROOT()), "/")
	defaultStaticFS   vfs.NameSpace = getNameSpace(mapfs.New(static.Files), "/")
	defaultDocFS      vfs.NameSpace = getNameSpace(defaultRootFS, "/doc")
	defaultLocalFS    vfs.NameSpace = getNameSpace(defaultRootFS, "/"+Default)
	defaultTranslater Translater    = new(localTranslater)
)

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
		defaultLocalFS = getNameSpace(defaultRootFS, "/"+Default)
	} else {
		if goRoot != "" && goRoot != runtime.GOROOT() {
			defaultRootFS = getNameSpace(vfs.OS(goRoot), "/")
			defaultDocFS = getNameSpace(defaultRootFS, "/doc")
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
