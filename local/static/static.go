// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

import (
	"github.com/chai2010/golangdoc/godoc/static"
	"github.com/chai2010/golangdoc/godoc/vfs/mapfs"
	"github.com/chai2010/golangdoc/local"
)

func Files(lang ...string) map[string]string {
	if len(lang) > 0 {
		if files, ok := StaticFilesTable[lang[0]]; ok && files != nil {
			return files
		}
	}
	return static.Files
}

func init() {
	for lang, files := range StaticFilesTable {
		local.RegisterStaticFS(lang, mapfs.New(files))
	}
}
