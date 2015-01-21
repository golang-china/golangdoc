// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"go/doc"

	"github.com/chai2010/golangdoc/godoc/vfs"
)

type nilTranslater struct{}

func (p *nilTranslater) Static(lang string) vfs.FileSystem {
	return nil
}

func (p *nilTranslater) Document(lang string) vfs.FileSystem {
	return nil
}

func (p *nilTranslater) Package(lang, importPath string, pkg ...*doc.Package) *doc.Package {
	return nil
}
