// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafe_zh_CN

import (
	"go/doc"
)

var (
	OriginPackage       *Package     = originPackage
	TranslatePackage    *Package     = translatePackage
	OriginDocPackage    *doc.Package = ToDocPackage(originPackage)
	TranslateDocPackage *doc.Package = ToDocPackage(translatePackage)
)

type Package struct {
	Name       string
	ImportPath string
	Doc        string
	Notes      []NoteValue
	Consts     []ConstValue
	Types      []TypeValue
	Vars       []VarValue
	Funcs      []FuncValue
	Methods    []MethodValue
}

type NoteValue struct {
	Type string // BUG/TODO/ ...
	UID  string // uid found with the marker
	Body string // note body text
}

type ConstValue struct {
	Type  string
	Names []string
	Doc   string
}

type TypeValue struct {
	Name string
	Doc  string
}

type VarValue struct {
	Type  string
	Names []string
	Doc   string
}

type FuncValue struct {
	Type string // for type Constructor
	Name string
	Doc  string
}

type MethodValue struct {
	Type string
	Name string
	Doc  string
}

func ToDocPackage(p *Package) *doc.Package {
	return &doc.Package{
		Name:       p.Name,
		ImportPath: p.ImportPath,
		Doc:        p.Doc,
		Notes:      p.makeNotes(),
		Consts:     p.makeConsts(),
		Types:      p.makeTypes(),
		Vars:       p.makeVars(),
		Funcs:      p.makeFuncs(),
	}
}

func (p *Package) makeNotes() map[string][]*doc.Note {
	for i := 0; i < len(p.Notes); i++ {
		if p.Notes[i].Type == "_" {
			p.Notes = append(p.Notes[:i], p.Notes[i+1:]...)
			i--
		}
	}
	if len(p.Notes) == 0 {
		return nil
	}
	notesMap := make(map[string][]*doc.Note)
	for i := 0; i < len(p.Notes); i++ {
		notes, _ := notesMap[p.Notes[i].Type]
		notes = append(notes, &doc.Note{
			UID:  p.Notes[i].UID,
			Body: p.Notes[i].Body,
		})
		notesMap[p.Notes[i].Type] = notes
	}
	return notesMap
}

func (p *Package) makeConsts() []*doc.Value {
	for i := 0; i < len(p.Consts); i++ {
		if len(p.Consts[i].Names) == 0 || p.Consts[i].Names[0] == "_" {
			p.Consts = append(p.Consts[:i], p.Consts[i+1:]...)
			i--
		}
	}
	if len(p.Consts) == 0 {
		return nil
	}
	consts := make([]*doc.Value, len(p.Consts))
	for i := 0; i < len(consts); i++ {
		consts[i] = &doc.Value{
			Names: append([]string(nil), p.Consts[i].Names...),
			Doc:   p.Consts[i].Doc,
		}
	}
	return consts
}

func (p *Package) makeTypes() []*doc.Type {
	for i := 0; i < len(p.Types); i++ {
		if p.Types[i].Name == "_" {
			p.Types = append(p.Types[:i], p.Types[i+1:]...)
			i--
		}
	}
	for i := 0; i < len(p.Methods); i++ {
		if p.Methods[i].Name == "_" {
			p.Methods = append(p.Methods[:i], p.Methods[i+1:]...)
			i--
		}
	}
	if len(p.Types) == 0 {
		return nil
	}
	types := make([]*doc.Type, len(p.Types))
	for i := 0; i < len(types); i++ {
		types[i] = &doc.Type{
			Name: p.Types[i].Name,
			Doc:  p.Types[i].Doc,

			Consts:  p.makeTypeConsts(p.Types[i].Name),
			Vars:    p.makeTypeVars(p.Types[i].Name),
			Funcs:   p.makeTypeFuncs(p.Types[i].Name),
			Methods: p.makeTypeMethods(p.Types[i].Name),
		}
	}
	return types
}

func (p *Package) makeFuncs() []*doc.Func {
	for i := 0; i < len(p.Funcs); i++ {
		if p.Funcs[i].Name == "_" {
			p.Funcs = append(p.Funcs[:i], p.Funcs[i+1:]...)
			i--
		}
	}
	if len(p.Funcs) == 0 {
		return nil
	}
	funcs := make([]*doc.Func, len(p.Funcs))
	for i := 0; i < len(funcs); i++ {
		funcs[i] = &doc.Func{
			Name: p.Funcs[i].Name,
			Doc:  p.Funcs[i].Doc,
		}
	}
	return funcs
}

func (p *Package) makeTypeConsts(typeName string) []*doc.Value {
	var consts []*doc.Value
	for i := 0; i < len(p.Consts); i++ {
		if p.Consts[i].Type == typeName {
			consts = append(consts, &doc.Value{
				Names: append([]string(nil), p.Consts[i].Names...),
				Doc:   p.Consts[i].Doc,
			})
		}
	}
	return consts
}

func (p *Package) makeTypeVars(typeName string) []*doc.Value {
	var vars []*doc.Value
	for i := 0; i < len(p.Vars); i++ {
		if p.Vars[i].Type == typeName {
			vars = append(vars, &doc.Value{
				Names: append([]string(nil), p.Vars[i].Names...),
				Doc:   p.Vars[i].Doc,
			})
		}
	}
	return vars
}

func (p *Package) makeTypeFuncs(typeName string) []*doc.Func {
	var funcs []*doc.Func
	for i := 0; i < len(p.Funcs); i++ {
		if p.Funcs[i].Type == typeName {
			funcs = append(funcs, &doc.Func{
				Name: p.Funcs[i].Name,
				Doc:  p.Funcs[i].Doc,
			})
		}
	}
	return funcs
}

func (p *Package) makeTypeMethods(typeName string) []*doc.Func {
	for i := 0; i < len(p.Methods); i++ {
		if p.Methods[i].Name == "_" {
			p.Methods = append(p.Methods[:i], p.Methods[i+1:]...)
			i--
		}
	}
	if len(p.Methods) == 0 {
		return nil
	}
	var methods []*doc.Func
	for i := 0; i < len(p.Methods); i++ {
		if p.Methods[i].Type == typeName {
			methods = append(methods, &doc.Func{
				Name: p.Methods[i].Name,
				Doc:  p.Methods[i].Doc,
			})
		}
	}
	return methods
}

func (p *Package) makeVars() []*doc.Value {
	for i := 0; i < len(p.Vars); i++ {
		if len(p.Vars[i].Names) == 0 || p.Vars[i].Names[0] == "_" {
			p.Vars = append(p.Vars[:i], p.Vars[i+1:]...)
			i--
		}
	}
	if len(p.Vars) == 0 {
		return nil
	}
	vars := make([]*doc.Value, len(p.Vars))
	for i := 0; i < len(vars); i++ {
		vars[i] = &doc.Value{
			Names: append([]string(nil), p.Vars[i].Names...),
			Doc:   p.Vars[i].Doc,
		}
	}
	return vars
}
