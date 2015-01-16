// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgdoc

import (
	"fmt"
	"go/doc"
)

func Register(local string, pkgDoc *doc.Package) error {
	return fmt.Errorf("pkgdoc: TODO")
}

func Getdoc(local string, pkgDoc *doc.Package) *doc.Package {
	return pkgDoc
}
