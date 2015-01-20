// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafe_zh_CN

var originPackage = &Package{
	Name:       "unsafe",
	ImportPath: "unsafe",
	Doc: `
	Package unsafe contains operations that step around the type safety of Go programs.

	Packages that import unsafe may be non-portable and are not protected by the
	Go 1 compatibility guidelines.
`,
	Notes:   originPackageNotes,
	Consts:  originPackageConsts,
	Types:   originPackageTypes,
	Vars:    originPackageVars,
	Funcs:   originPackageFuncs,
	Methods: originPackageMethods,
}

var originPackageNotes = []NoteValue{
	NoteValue{
		Type: "_",
		UID:  "_",
		Body: "_",
	},
}

var originPackageConsts = []ConstValue{
	ConstValue{
		Type:  "_",
		Names: []string{"_"},
		Doc:   "_",
	},
}

var originPackageTypes = []TypeValue{
	TypeValue{
		Name: "_",
		Doc:  "_",
	},
	TypeValue{
		Name: "ArbitraryType",
		Doc: `
// ArbitraryType is here for the purposes of documentation only and is not actually
// part of the unsafe package.  It represents the type of an arbitrary Go expression.
`,
	},
	TypeValue{
		Name: "Pointer",
		Doc: `
// Pointer represents a pointer to an arbitrary type.  There are four special operations
// available for type Pointer that are not available for other types.
//	1) A pointer value of any type can be converted to a Pointer.
//	2) A Pointer can be converted to a pointer value of any type.
//	3) A uintptr can be converted to a Pointer.
//	4) A Pointer can be converted to a uintptr.
// Pointer therefore allows a program to defeat the type system and read and write
// arbitrary memory. It should be used with extreme care.
`,
	},
}

var originPackageVars = []VarValue{
	VarValue{
		Type:  "_",
		Names: []string{"_"},
		Doc:   "_",
	},
}

var originPackageFuncs = []FuncValue{
	FuncValue{
		Type: "_",
		Name: "_",
		Doc:  "_",
	},
	FuncValue{
		Type: "uintptr",
		Name: "Sizeof",
		Doc: `
// Sizeof returns the size in bytes occupied by the value v.  The size is that of the
// "top level" of the value only.  For instance, if v is a slice, it returns the size of
// the slice descriptor, not the size of the memory referenced by the slice.
`,
	},
	FuncValue{
		Type: "uintptr",
		Name: "Offsetof",
		Doc: `
// Offsetof returns the offset within the struct of the field represented by v,
// which must be of the form structValue.field.  In other words, it returns the
// number of bytes between the start of the struct and the start of the field.
`,
	},
	FuncValue{
		Type: "uintptr",
		Name: "Alignof",
		Doc: `
// Alignof returns the alignment of the value v.  It is the maximum value m such
// that the address of a variable with the type of v will always be zero mod m.
// If v is of the form structValue.field, it returns the alignment of field f within struct object obj.
`,
	},
}

var originPackageMethods = []MethodValue{
	MethodValue{
		Type: "_",
		Name: "_",
		Doc:  "_",
	},
}
