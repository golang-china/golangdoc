// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run generate.go

package pkgdoc

import (
	"go/doc"
)

func DocPackage(p *doc.Package) *doc.Package {
	return p
}
