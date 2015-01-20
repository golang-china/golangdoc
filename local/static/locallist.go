// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

import (
	staticFiles_default "github.com/chai2010/golangdoc/local/static/default"
	staticFiles_zh_CN "github.com/chai2010/golangdoc/local/static/zh_CN"
)

var StaticFilesTable = map[string]map[string]string{
	"default": staticFiles_default.Files,
	"zh_CN":   staticFiles_zh_CN.Files,
}
