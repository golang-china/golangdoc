// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// local 包用于提供 godoc 的翻译支持.
package local

// 翻译文件所在的默认目录.
const (
	DefaultDir = "translations"     // $(RootFS)/translations
	DefaultEnv = "GODOC_LOCAL_ROOT" // dir list
)

// DocumentFS 函数返回 doc 目录对应的文件系统.
func DocumentFS(lang string) vfs.FileSystem

// DocumentFS 函数初始化本地翻译资源的环境, 支持zip格式.
func Init(goroot, gopath, translations string)

// Package 函数返回翻译后的包结构体.
func Package(lang, importPath string, pkg *doc.Package) *doc.Package
