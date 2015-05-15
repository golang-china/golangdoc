// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"archive/zip"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

const (
	__pkg__  = "__pkg__"
	__name__ = "__name__"
	__doc__  = "__doc__"
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

func openZipFS(filename string) vfs.NameSpace {
	rc, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatalf("local: zip.OpenReader(%q) failed: %s\n", filename, err)
	}
	fs := getNameSpace(zipfs.New(rc, filename), "/")

	rootdir := "/" + strings.TrimSuffix(
		filepath.Base(filename),
		filepath.Ext(filename),
	)
	if _, err = fs.Lstat(rootdir); err == nil {
		fs = getNameSpace(fs, rootdir)
	}
	return fs
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

func fsFileExists(fs vfs.NameSpace, name string) bool {
	if fi, err := fs.Stat(name); err != nil || fi.IsDir() {
		return false
	}
	return true
}

func mapKey(lang, importPath, id string) string {
	return fmt.Sprintf("%s.%s@%s", importPath, id, lang)
}

func methodId(typeName, methodName string) string {
	return typeName + "." + methodName
}

func parsePkgDocPackage(fs vfs.NameSpace, lang, importPath string) *doc.Package {
	if lang == "" || importPath == "" || importPath[0] == '/' {
		return nil
	}
	docCode := loadPkgDocCode(fs, lang, importPath)
	if docCode == nil {
		return nil
	}

	// parse doc
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, importPath, docCode, parser.ParseComments)
	if err != nil {
		log.Printf("local.localTranslater.ParseDocPackage: err = %v\n", err)
		return nil
	}
	astPkg, _ := ast.NewPackage(fset,
		map[string]*ast.File{importPath: astFile},
		nil,
		nil,
	)
	docPkg := doc.New(astPkg, importPath, doc.AllDecls)
	return docPkg
}

func loadPkgDocCode(fs vfs.NameSpace, lang, importPath string) []byte {
	// {FS}:/src/importPath/doc_$(lang)_GOOS_GOARCH.go
	// {FS}:/src/importPath/doc_$(lang)_GOARCH.go
	// {FS}:/src/importPath/doc_$(lang)_GOOS.go
	// {FS}:/src/importPath/doc_$(lang).go
	filenames := []string{
		fmt.Sprintf("/src/%s/doc_%s_%s_%s.go", importPath, lang, defaultGodocGoos, defaultGodocGoarch),
		fmt.Sprintf("/src/%s/doc_%s_%s.go", importPath, lang, defaultGodocGoarch),
		fmt.Sprintf("/src/%s/doc_%s_%s.go", importPath, lang, defaultGodocGoos),
		fmt.Sprintf("/src/%s/doc_%s.go", importPath, lang),
	}

	for i := 0; i < len(filenames); i++ {
		if fsFileExists(fs, filenames[i]) {
			docCode, _ := vfs.ReadFile(fs, filenames[i])
			if docCode != nil {
				return docCode
			}
		}
	}

	return nil
}
