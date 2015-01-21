// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package local

import (
	"github.com/chai2010/golangdoc/godoc/vfs"
)

type documentTranslater struct {
	nilTranslater
}

func (p *documentTranslater) Document(lang string) vfs.FileSystem {
	return nil
}
