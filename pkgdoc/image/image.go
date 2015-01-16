// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// AUTO GENERATED, DONOT EDIT!!!

package image

import (
	"go/doc"

	"github.com/chai2010/golangdoc/pkgdoc"
)

func init() {
	pkgdoc.Register("default", &doc.Package{
		Name:       pkgName,
		ImportPath: pkgImportPath,
		Doc:        pkgDoc,
		Notes:      pkgNotes,
		Bugs:       pkgBugs,
		Consts:     pkgConsts,
		Types:      pkgTypes,
		Vars:       pkgVars,
		Funcs:      pkgFuncs,
	})
}

var pkgName = "image"

var pkgImportPath = "image"

var pkgDoc = `
    Package image implements a basic 2-D image library.

    The fundamental interface is called Image. An Image contains colors,
    which are described in the image/color package.

    Values of the Image interface are created either by calling functions
    such as NewRGBA and NewPaletted, or by calling Decode on an io.Reader
    containing image data in a format such as GIF, JPEG or PNG. Decoding any
    particular image format requires the prior registration of a decoder
    function. Registration is typically automatic as a side effect of
    initializing that format's package so that, to decode a PNG image, it
    suffices to have

	import _ "image/png"

    in a program's main package. The _ means to import a package purely for
    its initialization side effects.

    See "The Go image package" for more details:
    http://golang.org/doc/articles/image_package.html
`

var pkgNotes = map[string][]*doc.Note{
	"TODO": []*doc.Note{
		&doc.Note{
			UID:  "chai2010",
			Body: "example",
		},
	},
}

var pkgBugs = []string{}

var pkgConsts = []*doc.Value{
	&doc.Value{
		Doc:   "",
		Names: []string{},
	},
}

var pkgTypes = []*doc.Type{
	&doc.Type{
		Doc:  "",
		Name: "",

		Consts: []*doc.Value{
			&doc.Value{
				Doc:   "",
				Names: []string{},
			},
		},
		Vars: []*doc.Value{
			&doc.Value{
				Doc:   "",
				Names: []string{},
			},
		},
		Funcs: []*doc.Func{
			&doc.Func{
				Doc:  "",
				Name: "",
			},
		},
		Methods: []*doc.Func{
			&doc.Func{
				Doc:  "",
				Name: "",
			},
		},
	},
}

var pkgVars = []*doc.Value{
	&doc.Value{
		Doc:   "",
		Names: []string{},
	},
}

var pkgFuncs = []*doc.Func{
	&doc.Func{
		Doc:  "",
		Name: "",
	},
}
