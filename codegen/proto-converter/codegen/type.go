package codegen

import (
	"github.com/dave/jennifer/jen"
)

type Type interface {
	// PbType is the type in the protoc-generated code.
	PbType() jen.Code
	// MirType is the type in the Mir-generated code.
	// If the same type is used in protoc-generated structs and Mir-generated structs, PbType() should be equal
	// (w.r.t. "==") to MirType().
	MirType() jen.Code
	// ToMir converts an object of the type in protoc representation to its Mir representation.
	ToMir(code jen.Code) jen.Code
	// ToPb converts an object of the type in Mir representation to its protoc representation.
	ToPb(code jen.Code) jen.Code
}

// Same is used when the same type is used by protoc-generated code and Mir-generated code.
type Same struct {
	tp jen.Code
}

func (t Same) PbType() jen.Code {
	return t.tp
}

func (t Same) MirType() jen.Code {
	return t.tp
}

func (t Same) ToMir(code jen.Code) jen.Code {
	return code
}

func (t Same) ToPb(code jen.Code) jen.Code {
	return code
}

// Castable is used when the types used by protoc-generated code and
// Mir-generated code can be directly cast to one another.
type Castable struct {
	pbType  jen.Code
	mirType jen.Code
}

func (t Castable) PbType() jen.Code {
	return t.pbType
}

func (t Castable) MirType() jen.Code {
	return t.mirType
}

func (t Castable) ToMir(code jen.Code) jen.Code {
	return jen.Parens(t.mirType).Call(code)
}

func (t Castable) ToPb(code jen.Code) jen.Code {
	return jen.Parens(t.pbType).Call(code)
}
