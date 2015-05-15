// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"golang.org/x/tools/godoc/vfs"
)

func getNameSpace(fs vfs.FileSystem, ns string) vfs.NameSpace {
	newns := make(vfs.NameSpace)
	if ns != "" {
		newns.Bind("/", fs, ns, vfs.BindReplace)
	} else {
		newns.Bind("/", fs, "/", vfs.BindReplace)
	}
	return newns
}
