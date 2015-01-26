// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"text/template"
)

var testDocTemplate = template.Must(
	template.New("doc").Funcs(template.FuncMap{
		"comment_text": comment_textFunc,
		"node":         nodeFunc,
	}).Parse(
		tmplPackageText,
	),
)

func TestDocgen(t *testing.T) {
	//
}
