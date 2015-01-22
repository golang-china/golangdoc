// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// local 包用于提供 godoc 的翻译支持.
package local

// 翻译文件所在的默认目录.
const (
	Default = "translations" // $(RootFS)/translations
)

// DocumentFS 函数返回 doc 目录对应的文件系统.
func DocumentFS(lang string) vfs.FileSystem

// DocumentFS 函数初始化本地翻译资源的环境.
func Init(goRoot, goZipFile, goTemplateDir string)

// Package 函数返回翻译后的包结构体.
func Package(lang, importPath string, pkg ...*doc.Package) *doc.Package

// RegisterDocumentFS 函数注册 doc 目录的翻译版本.
func RegisterDocumentFS(lang string, docFiles vfs.FileSystem)

// RegisterPackage 函数注册包对应的翻译信息.
func RegisterPackage(lang string, pkg *doc.Package)

// RegisterStaticFS 函数注册静态文件的翻译版本.
func RegisterStaticFS(lang string, staticFiles vfs.FileSystem)

// RegisterTranslater 函数注册翻译器.
func RegisterTranslater(tr Translater)

// RootFS 函数返回根目录文件系统.
func RootFS() vfs.FileSystem

// StaticFS 函数返回静态文件的文件系统.
func StaticFS(lang string) vfs.FileSystem

// Translater 为翻译器接口类型.
type Translater interface {
	Static(lang string) vfs.FileSystem
	Document(lang string) vfs.FileSystem
	Package(lang, importPath string, pkg ...*doc.Package) *doc.Package
}
