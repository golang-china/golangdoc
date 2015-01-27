// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Generate Go pakckage doc for translate.
//
// Usage:
//	docgen package lang... [-GOOS=...] [-GOARCH=...]
//	docgen -h
//
// Example:
//	docgen builtin zh_CN
//	docgen unsafe  zh_CN
//	docgen unsafe  zh_CN
//	docgen syscall zh_CN -GOOS=windows                     # for windows
//	docgen syscall zh_CN -GOOS=windows -GOARCH=amd64       # for windows/amd64
//	docgen syscall zh_CN                                   # for non windows
//	docgen std     zh_CN                                   # all standard packages
//	docgen ./...   zh_CN                                   # all sub packages
//
// Output:
//	translations/src/builtin/doc_zh_CN.go
//	translations/src/unsafe/doc_zh_CN.go
//	translations/src/unsafe/doc_zh_CN.go
//	translations/src/syscall/doc_zh_CN_windows.go          # for windows
//	translations/src/syscall/doc_zh_CN_windows_amd64.go    # for windows/amd64
//	translations/src/syscall/doc_zh_CN.go                  # for non windows
//	translations/src/*/doc_zh_CN.go                        # all standard packages
//	translations/src/*/doc_zh_CN.go                        # all sub packages
//
// Help:
//	docgen -h
//
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/chai2010/golangdoc/local"
)

const usage = `
Usage: docgen package lang... [-GOOS=...] [-GOARCH=...]
  docgen -h

Example:
  docgen builtin zh_CN
  docgen unsafe  zh_CN
  docgen unsafe  zh_CN 
  docgen syscall zh_CN -GOOS=windows                     # for windows
  docgen syscall zh_CN -GOOS=windows -GOARCH=amd64       # for windows/amd64
  docgen syscall zh_CN                                   # for non windows
  docgen std     zh_CN                                   # all standard packages
  docgen ./...   zh_CN                                   # all sub packages

Output:
  translations/src/builtin/doc_zh_CN.go
  translations/src/unsafe/doc_zh_CN.go
  translations/src/unsafe/doc_zh_CN.go
  translations/src/syscall/doc_zh_CN_windows.go          # for windows
  translations/src/syscall/doc_zh_CN_windows_amd64.go    # for windows/amd64
  translations/src/syscall/doc_zh_CN.go                  # for non windows
  translations/src/*/doc_zh_CN.go                        # all standard packages
  translations/src/*/doc_zh_CN.go                        # all sub packages

Help:
  docgen -h
    
Report bugs to <chaishushan{AT}gmail.com>.
`

var (
	flagGOOS       = ""
	flagGOARCH     = ""
	cmdArgPackages = []string(nil)
	cmdArgLangs    = []string(nil)
)

func main() {
	parseCmdArgs()
	for i := 0; i < len(cmdArgPackages); i++ {
		for _, lang := range cmdArgLangs {
			if importPath, err := docgen(cmdArgPackages[i], lang); err != nil {
				log.Fatalf("gen %s failed, err = %v", docFilename(importPath, lang), err)
			} else {
				fmt.Printf("gen %s ok\n", docFilename(importPath, lang))
			}
		}
	}
	fmt.Println("Done")
}

func parseCmdArgs() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, usage[1:len(usage)-1])
		os.Exit(0)
	}
	var args []string
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-h" || os.Args[i] == "-help" {
			fmt.Fprintln(os.Stderr, usage[1:len(usage)-1])
			os.Exit(0)
		}
		if strings.HasPrefix(os.Args[i], "-GOOS=") {
			flagGOOS = os.Args[i][len("-GOOS="):]
			continue
		}
		if strings.HasPrefix(os.Args[i], "-GOARCH=") {
			flagGOARCH = os.Args[i][len("-GOARCH="):]
			continue
		}
		args = append(args, os.Args[i])
	}
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, usage[1:len(usage)-1])
		os.Exit(1)
	}

	cmdArgPackages = listPackages(args[0])
	cmdArgLangs = args[1:]
}

func listPackages(name string) (pkgs []string) {
	listOut, err := exec.Command(`go`, `list`, name).Output()
	if err != nil {
		log.Fatalf("listPackages: err = %v", err)
	}
	for _, line := range strings.Split(string(listOut), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			pkgs = append(pkgs, s)
		}
	}
	if name == "std" {
		hasBuiltin := false
		for _, s := range pkgs {
			if s == "builtin" {
				hasBuiltin = true
				break
			}
		}
		if !hasBuiltin {
			pkgs = append(pkgs, "builtin")
		}
	}
	sort.Strings(pkgs)
	return
}

func docgen(name, lang string) (importPath string, err error) {
	info, err := ParsePackageInfo(name, lang)
	if err != nil {
		return
	}
	importPath = info.PDoc.ImportPath

	filename := docFilename(importPath, lang)
	os.MkdirAll(path.Dir(filename), 0666)

	data, err := format.Source(info.Bytes())
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return
	}
	return
}

func docFilename(importPath, lang string) string {
	const base = "translations/src"
	if importPath == "syscall" && flagGOOS == "" {
		flagGOOS = runtime.GOOS
	}
	var filename string
	switch {
	case flagGOOS != "" && flagGOARCH != "":
		filename = path.Join(base, importPath, fmt.Sprintf("doc_%s_%s_%s.go", lang, flagGOOS, flagGOARCH))
	case flagGOOS != "":
		filename = path.Join(base, importPath, fmt.Sprintf("doc_%s_%s.go", lang, flagGOOS))
	case flagGOARCH != "":
		filename = path.Join(base, importPath, fmt.Sprintf("doc_%s_%s.go", lang, flagGOARCH))
	default:
		filename = path.Join(base, importPath, fmt.Sprintf("doc_%s.go", lang))
	}
	return filename
}

type PackageInfo struct {
	Lang      string
	FSet      *token.FileSet
	PAst      *ast.Package
	PDoc      *doc.Package
	PDocLocal *doc.Package
	PDocMap   map[string]string
}

func ParsePackageInfo(name, lang string) (pkg *PackageInfo, err error) {
	type PkgInfo struct {
		Dir        string // directory containing package sources
		Name       string // package name
		ImportPath string // import path of package in dir
	}

	var pkgInfo PkgInfo
	listOut, err := exec.Command(`go`, `list`, `-json`, name).Output()
	if err != nil {
		return
	}
	err = json.Unmarshal(listOut, &pkgInfo)
	if err != nil {
		return
	}

	fset := token.NewFileSet()
	past, err := parser.ParseDir(fset, pkgInfo.Dir,
		func(fi os.FileInfo) bool {
			if strings.HasSuffix(fi.Name(), "_test.go") {
				return false
			}
			return true
		},
		parser.ParseComments,
	)
	if err != nil {
		return
	}

	var mode doc.Mode
	if pkgInfo.ImportPath == "builtin" {
		mode = doc.AllDecls
	}
	pdoc := doc.New(past[pkgInfo.Name], pkgInfo.ImportPath, mode)
	pdocLocal := local.Package(lang, pkgInfo.ImportPath)

	pkg = &PackageInfo{
		Lang:      lang,
		FSet:      fset,
		PAst:      past[pkgInfo.Name],
		PDoc:      pdoc,
		PDocLocal: pdocLocal,
		PDocMap:   make(map[string]string),
	}
	pkg.initDocTable("", pkg.PDoc)
	if pkg.PDocLocal != nil {
		pkg.initDocTable(lang, pkg.PDocLocal)
	}
	return
}

func (p *PackageInfo) mapKey(lang, importPath, id string) string {
	return fmt.Sprintf("%s.%s@%s", importPath, id, lang)
}

func (p *PackageInfo) methodId(typeName, methodName string) string {
	return typeName + "." + methodName
}

func (p *PackageInfo) initDocTable(lang string, pkg *doc.Package) {
	for _, v := range pkg.Consts {
		for _, id := range v.Names {
			p.PDocMap[p.mapKey(lang, pkg.ImportPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Types {
		p.PDocMap[p.mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc

		for _, x := range v.Consts {
			for _, id := range x.Names {
				p.PDocMap[p.mapKey(lang, pkg.ImportPath, id)] = x.Doc
			}
		}
		for _, x := range v.Vars {
			for _, id := range x.Names {
				p.PDocMap[p.mapKey(lang, pkg.ImportPath, id)] = x.Doc
			}
		}
		for _, x := range v.Funcs {
			p.PDocMap[p.mapKey(lang, pkg.ImportPath, x.Name)] = x.Doc
		}
		for _, x := range v.Methods {
			p.PDocMap[p.mapKey(lang, pkg.ImportPath, p.methodId(v.Name, x.Name))] = x.Doc
		}
	}
	for _, v := range pkg.Vars {
		for _, id := range v.Names {
			p.PDocMap[p.mapKey(lang, pkg.ImportPath, id)] = v.Doc
		}
	}
	for _, v := range pkg.Funcs {
		p.PDocMap[p.mapKey(lang, pkg.ImportPath, v.Name)] = v.Doc
	}
}

func (p *PackageInfo) Bytes() []byte {
	var docTemplate = template.Must(
		template.New("doc").Funcs(template.FuncMap{
			"comment_text": p.comment_textFunc,
			"node":         p.nodeFunc,
		}).Parse(
			tmplPackageText,
		),
	)

	var out bytes.Buffer
	if err := docTemplate.Execute(&out, p); err != nil {
		log.Fatal(fmt.Sprintf("PackageInfo.Bytes: err = %v", err))
	}
	return out.Bytes()
}

func (p *PackageInfo) comment_textFunc(id, comment, indent, preIndent string) string {
	localDoc := p.getLocalDoc(id)
	comment1 := p.comment_format(comment, indent, preIndent)
	if localDoc == "" || localDoc == comment {
		return comment1
	}
	comment2 := p.comment_format(localDoc, indent, preIndent)
	if comment1 == comment2 {
		return comment1
	}
	return comment1 + "\n\n" + comment2
}

func (p *PackageInfo) getLocalDoc(id string) string {
	if p.PDocLocal == nil {
		return ""
	}
	if id != "" {
		s, _ := p.PDocMap[p.mapKey(p.Lang, p.PDoc.ImportPath, id)]
		return s
	} else {
		return p.PDocLocal.Doc
	}
}

func (p *PackageInfo) comment_format(comment, indent, preIndent string) string {
	containsOnlySpace := func(buf []byte) bool {
		isNotSpace := func(r rune) bool { return !unicode.IsSpace(r) }
		return bytes.IndexFunc(buf, isNotSpace) == -1
	}
	var buf bytes.Buffer
	const punchCardWidth = 80
	ToText(&buf, comment, indent, preIndent, punchCardWidth-2*len(indent))
	if containsOnlySpace(buf.Bytes()) {
		return ""
	}
	lines := strings.Split(buf.String(), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" || lines[i][0] != '\t' {
			lines[i] = "// " + lines[i]
		} else {
			lines[i] = "//" + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}

func (p *PackageInfo) nodeFunc(node interface{}) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, p.FSet, node)
	if err != nil {
		log.Print(err)
	}
	return buf.String()
}

const tmplPackageText = `// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

{{with .PDoc}}{{/* template comments */}}{{/*

-------------------------------------------------------------------------------
-- PACKAGE DOCUMENTATION
-------------------------------------------------------------------------------

*/}}{{comment_text "" .Doc "" "\t"}}
package {{.Name}}
{{/*

-------------------------------------------------------------------------------
-- CONSTANTS
-------------------------------------------------------------------------------

*/}}{{with .Consts}}{{range .}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- VARIABLES
-------------------------------------------------------------------------------

*/}}{{with .Vars}}{{range .}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- FUNCTIONS
-------------------------------------------------------------------------------

*/}}{{with .Funcs}}{{range .}}
{{comment_text .Name .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES
-------------------------------------------------------------------------------

*/}}{{with .Types}}{{range .}}{{$typeName := .Name}}
{{comment_text .Name .Doc "" "\t"}}
{{node .Decl}}
{{/*

-------------------------------------------------------------------------------
-- TYPES.CONSTANTS
-------------------------------------------------------------------------------

*/}}{{if .Consts}}{{range .Consts}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.VARIABLES
-------------------------------------------------------------------------------

*/}}{{if .Vars}}{{range .Vars}}
{{comment_text (index .Names 0) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.FUNCTIONS
-------------------------------------------------------------------------------

*/}}{{if .Funcs}}{{range .Funcs}}
{{comment_text .Name .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
TYPES.METHODS
-------------------------------------------------------------------------------

*/}}{{if .Methods}}{{range .Methods}}
{{comment_text (printf "%s.%s" $typeName .Name) .Doc "" "\t"}}
{{node .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.END
-------------------------------------------------------------------------------

*/}}{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- END
-------------------------------------------------------------------------------

*/}}{{end}}
`
