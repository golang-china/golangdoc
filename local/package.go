// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"go/doc"
)

type packageTranslater struct {
	nilTranslater
}

func (p *packageTranslater) Package(lang, importPath string, pkg ...*doc.Package) *doc.Package {
	return nil
}
