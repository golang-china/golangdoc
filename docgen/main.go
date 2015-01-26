// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Generate Go pakckage doc for translate.
//
// Usage:
//	docgen importPath lang... [-GOOS=...] [-GOARCH=...]
//	docgen -h
//
// Example:
//	docgen builtin zh_CN
//	docgen unsafe  zh_CN
//	docgen unsafe  zh_CN zh_TW
//	docgen syscall zh_CN zh_TW -GOOS=windows               # for windows
//	docgen syscall zh_CN zh_TW -GOOS=windows -GOARCH=amd64 # for windows/amd64
//	docgen syscall zh_CN zh_TW                             # for non windows
//
// Output:
//	translations/src/builtin/doc_zh_CN.go
//	translations/src/unsafe/doc_zh_CN.go unsafe/doc_zh_TW.go
//	translations/src/unsafe/doc_zh_CN.go
//	translations/src/syscall/doc_zh_CN_windows.go          # for windows
//	translations/src/syscall/doc_zh_CN_windows_amd64.go    # for windows/amd64
//	translations/src/syscall/doc_zh_CN.go                  # for non windows
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
	"strings"
	"text/template"
	"unicode"

	"github.com/chai2010/golangdoc/local"
)

const usage = `
Usage: docgen importPath lang... [-GOOS=...] [-GOARCH=...]
  docgen -h

Example:
  docgen builtin zh_CN
  docgen unsafe  zh_CN
  docgen unsafe  zh_CN zh_TW
  docgen syscall zh_CN zh_TW -GOOS=windows               # for windows
  docgen syscall zh_CN zh_TW -GOOS=windows -GOARCH=amd64 # for windows/amd64
  docgen syscall zh_CN zh_TW                             # for non windows

Output:
  translations/src/builtin/doc_zh_CN.go
  translations/src/unsafe/doc_zh_CN.go unsafe/doc_zh_TW.go
  translations/src/unsafe/doc_zh_CN.go
  translations/src/syscall/doc_zh_CN_windows.go          # for windows
  translations/src/syscall/doc_zh_CN_windows_amd64.go    # for windows/amd64
  translations/src/syscall/doc_zh_CN.go                  # for non windows

Help:
  docgen -h
    
Report bugs to <chaishushan{AT}gmail.com>.
`

var (
	flagGOOS         = ""
	flagGOARCH       = ""
	cmdArgImportPath = ""
	cmdArgLangList   = []string(nil)
)

func main() {
	parseCmdArgs()
	for _, lang := range cmdArgLangList {
		if err := docgen(cmdArgImportPath, lang); err != nil {
			log.Fatalf("gen %s failed, err = %v", docFilename(cmdArgImportPath, lang), err)
		}
		fmt.Printf("gen %s ok\n", docFilename(cmdArgImportPath, lang))
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
	cmdArgImportPath = args[0]
	cmdArgLangList = args[1:]
}

func docgen(importPath, lang string) error {
	info, err := ParsePackageInfo(cmdArgImportPath, lang)
	if err != nil {
		return err
	}

	filename := docFilename(importPath, lang)
	os.MkdirAll(path.Dir(filename), 0666)

	data, err := format.Source([]byte(info.String()))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func docFilename(importPath, lang string) string {
	const base = "translations/src"
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
	FSet      *token.FileSet
	PAst      *ast.Package
	PDoc      *doc.Package
	PDocLocal *doc.Package
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
	pdoc := doc.New(past[pkgInfo.Name], pkgInfo.ImportPath, 0)
	pdocLocal := local.Package(lang, pkgInfo.ImportPath, pdoc)

	pkg = &PackageInfo{
		FSet:      fset,
		PAst:      past[pkgInfo.Name],
		PDoc:      pdoc,
		PDocLocal: pdocLocal,
	}
	return
}

func (p *PackageInfo) String() string {
	var docTemplate = template.Must(
		template.New("doc").Funcs(template.FuncMap{
			"comment_text": comment_textFunc,
			"node":         nodeFunc,
		}).Parse(
			tmplPackageText,
		),
	)

	var out bytes.Buffer
	if err := docTemplate.Execute(&out, p); err != nil {
		return fmt.Sprintf("err = %v", err)
	}
	return out.String()
}

func comment_textFunc(comment, indent, preIndent string) string {
	containsOnlySpace := func(buf []byte) bool {
		isNotSpace := func(r rune) bool { return !unicode.IsSpace(r) }
		return bytes.IndexFunc(buf, isNotSpace) == -1
	}
	var buf bytes.Buffer
	const punchCardWidth = 80
	doc.ToText(&buf, comment, indent, preIndent, punchCardWidth-2*len(indent))
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

func nodeFunc(info *PackageInfo, node interface{}) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, info.FSet, node)
	if err != nil {
		log.Print(err)
	}
	return buf.String()
}

const tmplPackageText = `// Copyright The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

{{with .PDocLocal}}{{/* template comments */}}{{/*

-------------------------------------------------------------------------------
-- PACKAGE DOCUMENTATION
-------------------------------------------------------------------------------

*/}}{{comment_text .Doc "" "\t"}}
package {{.Name}}
{{/*

-------------------------------------------------------------------------------
-- CONSTANTS
-------------------------------------------------------------------------------

*/}}{{with .Consts}}{{range .}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- VARIABLES
-------------------------------------------------------------------------------

*/}}{{with .Vars}}{{range .}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- FUNCTIONS
-------------------------------------------------------------------------------

*/}}{{with .Funcs}}{{range .}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES
-------------------------------------------------------------------------------

*/}}{{with .Types}}{{range .}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
{{/*

-------------------------------------------------------------------------------
-- TYPES.CONSTANTS
-------------------------------------------------------------------------------

*/}}{{if .Consts}}{{range .Consts}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.VARIABLES
-------------------------------------------------------------------------------

*/}}{{if .Vars}}{{range .Vars}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
-- TYPES.FUNCTIONS
-------------------------------------------------------------------------------

*/}}{{if .Funcs}}{{range .Funcs}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
{{end}}{{end}}{{/*

-------------------------------------------------------------------------------
TYPES.METHODS
-------------------------------------------------------------------------------

*/}}{{if .Methods}}{{range .Methods}}
{{comment_text .Doc "" "\t"}}
{{node $ .Decl}}
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
