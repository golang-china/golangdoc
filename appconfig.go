// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine

package main

const (
	// .zip filename
	flagZipFilename = "goroot.zip"

	// goroot directory in .zip file
	flagZipGoroot = "goroot"

	// glob pattern describing search index files
	// (if empty, the index is built at run-time)
	flagIndexFilenames = ""
)

var flagLang = func() *string {
	v := "zh_CN"
	return &v
}()
