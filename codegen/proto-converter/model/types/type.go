package types

import (
	"github.com/dave/jennifer/jen"
)

type Type interface {
	// Same returns whether PbType() and MirType() represent the same type.
	Same() bool
	// PbType is the type in the protoc-generated code.
	PbType() *jen.Statement
	// MirType is the type in the Mir-generated code.
	MirType() *jen.Statement
	// ToMir converts an object of the type in protoc representation to its Mir representation.
	ToMir(code jen.Code) *jen.Statement
	// ToPb converts an object of the type in Mir representation to its protoc representation.
	ToPb(code jen.Code) *jen.Statement
}

// Same is used when the same type is used by protoc-generated code and Mir-generated code.
type Same struct {
	tp jen.Code
}

func (t Same) Same() bool {
	return true
}

func (t Same) PbType() *jen.Statement {
	return jen.Add(t.tp)
}

func (t Same) MirType() *jen.Statement {
	return jen.Add(t.tp)
}

func (t Same) ToMir(code jen.Code) *jen.Statement {
	return jen.Add(code)
}

func (t Same) ToPb(code jen.Code) *jen.Statement {
	return jen.Add(code)
}

// Castable is used when the types used by protoc-generated code and
// Mir-generated code can be directly cast to one another.
type Castable struct {
	pbType  jen.Code
	mirType jen.Code
}

func (t Castable) Same() bool {
	return false
}

func (t Castable) PbType() *jen.Statement {
	return jen.Add(t.pbType)
}

func (t Castable) MirType() *jen.Statement {
	return jen.Add(t.mirType)
}

func (t Castable) ToMir(code jen.Code) *jen.Statement {
	return jen.Parens(t.mirType).Call(code)
}

func (t Castable) ToPb(code jen.Code) *jen.Statement {
	return jen.Parens(t.pbType).Call(code)
}