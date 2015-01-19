// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

import (
	defaultStatic "github.com/chai2010/golangdoc/local/static/default"
	zh_CNStatic "github.com/chai2010/golangdoc/local/static/zh_CN"
)

var StaticFilesTable = map[string]map[string]string{
	"default": defaultStatic.Files,
	"zh_CN":   zh_CNStatic.Files,
}
