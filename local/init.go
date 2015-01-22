// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"archive/zip"
	"log"
	"runtime"

	"github.com/chai2010/golangdoc/godoc/static"
	"github.com/chai2010/golangdoc/godoc/vfs"
	"github.com/chai2010/golangdoc/godoc/vfs/mapfs"
	"github.com/chai2010/golangdoc/godoc/vfs/zipfs"
)

const (
	Default = "translations" // $(RootFS)/translations
)

var (
	defaultRootFS     vfs.FileSystem = vfs.OS(runtime.GOROOT())
	defaultStaticFS   vfs.FileSystem = mapfs.New(static.Files)
	defaultDocFS      vfs.FileSystem = getNameSpace(defaultRootFS, "/doc")
	defaultLocalFS    vfs.FileSystem = getNameSpace(defaultRootFS, "/"+Default)
	defaultTranslater Translater     = new(localTranslater)
)

func Init(goRoot, goZipFile, goTemplateDir string) {
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
			defaultRootFS = vfs.OS(goRoot)
			defaultDocFS = getNameSpace(defaultRootFS, "/doc")
			defaultLocalFS = getNameSpace(defaultRootFS, "/"+Default)
		}
	}

	if goTemplateDir != "" {
		defaultStaticFS = vfs.OS(goTemplateDir)
	}
}

func getNameSpace(fs vfs.FileSystem, ns string) vfs.FileSystem {
	if ns != "" {
		subfs := make(vfs.NameSpace)
		subfs.Bind("/", fs, ns, vfs.BindReplace)
		return subfs
	}
	return fs
}
