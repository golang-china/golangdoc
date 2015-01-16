// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package example is a example for pkgdoc.

Example:
	func main() {
		println("Hello, 世界")
	}

Please report bugs to <chaishushan{AT}gmail.com>.
Thanks!
*/
package example

// TODO(chai2010): Add more examples.

// BUG(chai2010): Just a bug example.

// const1 doc.
const (
	CONST1_A     = 'A' // doc CONST1_A?
	CONST1_B     = 'B'
	CONST1_Byte  = Byte('A')
	CONST1_MyInt = MyInt('A')
)

// const2 doc.
const (
	CONST2_A     = 'A'
	CONST2_B     = 'B'
	CONST2_Byte  = Byte('A')
	CONST2_MyInt = MyInt('A')
)

const CONST3_NO_DOC = "xxx"

// Byte type doc.
type Byte uint8

func PrintByte(a Byte) {
	//
}

// MyInt type doc.
type MyInt int

func PrintMyInt(a MyInt) {
	//
}

func (p MyInt) Print() {
	//
}
