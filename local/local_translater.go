// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"

	"golang.org/x/tools/godoc/vfs"
)

type localTranslater struct{}

func (p *localTranslater) Static(lang string) vfs.FileSystem {
	if lang == "" {
		return defaultStaticFS
	}
	return p.NameSpace("/static/" + lang)
}

func (p *localTranslater) Document(lang string) vfs.FileSystem {
	if lang == "" {
		return defaultDocFS
	}
	return p.NameSpace("/doc/" + lang)
}

func (p *localTranslater) Package(lang, importPath string, pkg ...*doc.Package) *doc.Package {
	if lang == "" {
		if len(pkg) > 0 {
			return pkg[0]
		} else {
			return nil
		}
	}

	// try parse and register new pkg doc
	localPkg := p.ParseDocPackage(lang, importPath)
	if localPkg == nil {
		return nil
	}
	RegisterPackage(lang, localPkg)

	// retry Package func
	return Package(lang, importPath, pkg...)
}

func (p *localTranslater) Blog(lang string) vfs.FileSystem {
	if lang == "" {
		return defaultBlogFS
	}
	return p.NameSpace("/blog/" + lang)
}

func (p *localTranslater) ParseDocPackage(lang, importPath string) *doc.Package {
	if lang == "" || importPath == "" || importPath[0] == '/' {
		return nil
	}
	docCode := p.loadDocCode(lang, importPath)
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

func (p *localTranslater) NameSpace(ns string) vfs.FileSystem {
	if ns != "" {
		if fi, err := defaultLocalFS.Stat(ns); err != nil || !fi.IsDir() {
			return nil
		}
		subfs := make(vfs.NameSpace)
		subfs.Bind("/", defaultLocalFS, ns, vfs.BindReplace)
		return subfs
	}
	return defaultLocalFS
}

func (p *localTranslater) loadDocCode(lang, importPath string) []byte {
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
		// $(GOROOT)/translates/
		if p.fileExists(defaultLocalFS, filenames[i]) {
			docCode, _ := vfs.ReadFile(defaultLocalFS, filenames[i])
			if docCode != nil {
				return docCode
			}
		}

		// $(GOROOT)/
		if p.fileExists(defaultRootFS, filenames[i]) {
			docCode, _ := vfs.ReadFile(defaultRootFS, filenames[i])
			if docCode != nil {
				return docCode
			}
		}
	}

	return nil
}

func (p *localTranslater) fileExists(fs vfs.NameSpace, name string) bool {
	if fi, err := fs.Stat(name); err != nil || fi.IsDir() {
		return false
	}
	return true
}
