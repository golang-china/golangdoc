// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafe_zh_CN

var translatePackage = &Package{
	Name:       "unsafe",
	ImportPath: "unsafe",
	Doc: `
	unsafe 包含有关于Go程序类型安全的所有操作.
`,
	Notes:   translatePackageNotes,
	Consts:  translatePackageConsts,
	Types:   translatePackageTypes,
	Vars:    translatePackageVars,
	Funcs:   translatePackageFuncs,
	Methods: translatePackageMethods,
}

var translatePackageNotes = []NoteValue{
	NoteValue{
		Type: "_",
		UID:  "_",
		Body: "_",
	},
}

var translatePackageConsts = []ConstValue{
	ConstValue{
		Type:  "_",
		Names: []string{"_"},
		Doc:   "_",
	},
}

var translatePackageTypes = []TypeValue{
	TypeValue{
		Name: "_",
		Doc:  "_",
	},
	TypeValue{
		Name: "ArbitraryType",
		Doc: `
// ArbitraryType 在此处只用作文档目的，它实际上并不是 unsafe 包的一部分。
// 它代表任意一个Go表达式的类型。
`,
	},
	TypeValue{
		Name: "Pointer",
		Doc: `
// Pointer 代表一个指向任意类型的指针。
// 有三种特殊的操作可用于类型指针而不能用于其它类型。
//	1) 任意类型的指针值均可转换为 Pointer。
//	2) Pointer 均可转换为任意类型的指针值。
//	3) uintptr 均可转换为 Pointer。
//	4) Pointer 均可转换为 uintptr。
// 因此 Pointer 允许程序击溃类型系统并读写任意内存。它应当被用得非常小心。
`,
	},
}

var translatePackageVars = []VarValue{
	VarValue{
		Type:  "_",
		Names: []string{"_"},
		Doc:   "_",
	},
}

var translatePackageFuncs = []FuncValue{
	FuncValue{
		Type: "_",
		Name: "_",
		Doc:  "_",
	},
	FuncValue{
		Type: "uintptr",
		Name: "Sizeof",
		Doc: `
// Sizeof 返回被值 v 所占用的字节大小。
// 该大小只是最“顶级”的值。例如，若 v 是一个切片，它会返回该切片描述符的大小，
// 而非该切片引用的内存大小。
`,
	},
	FuncValue{
		Type: "uintptr",
		Name: "Offsetof",
		Doc: `
// Offsetof 返回由 v 所代表的结构中字段的偏移，它必须为 structValue.field 的形式。
// 换句话说，它返回该结构起始处与该字段起始处之间的字节数。
`,
	},
	FuncValue{
		Type: "uintptr",
		Name: "Alignof",
		Doc: `
// Alignof 返回 v 值的对齐方式。
// 其返回值 m 满足变量 v 的类型地址与 m 取模为 0 的最大值。若 v 是 structValue.field
// 的形式，它会返回字段 f 在其相应结构对象 obj 中的对齐方式。
`,
	},
}

var translatePackageMethods = []MethodValue{
	MethodValue{
		Type: "_",
		Name: "_",
		Doc:  "_",
	},
}
