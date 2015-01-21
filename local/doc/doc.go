// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package doc

import (
	"os"
	"path"
	"runtime"

	"github.com/chai2010/golangdoc/godoc/vfs"
	"github.com/chai2010/golangdoc/local"
)

func init() {
	zhDocPath := path.Join(runtime.GOROOT(), local.BaseName, "zh_CN", "doc")
	if isDirExist(zhDocPath) {
		local.RegisterDocFS("zh_CN", vfs.OS(zhDocPath))
	}
}

func isDirExist(name string) bool {
	if fi, err := os.Stat(name); err != nil || !fi.IsDir() {
		return false
	}
	return true
}
