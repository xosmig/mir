//package codegen
//
//import (
//	"github.com/dave/jennifer/jen"
//)
//
//type Type interface {
//	// PbType is the type in the protoc-generated code.
//	PbType() jen.Code
//	// MirType is the type in the Mir-generated code.
//	MirType() jen.Code
//	// FromPbType converts an object of the type in protoc representation to its Mir representation.
//	FromPbType(code jen.Code) jen.Code
//	// ToPbType converts an object of the type in Mir representation to its protoc representation.
//	ToPbType(code jen.Code) jen.Code
//}
//
//type Slice struct {
//	Underlying Type
//}
//
//// Implementation of Type interface for Slice.
//
//func (s *Slice) PbType() jen.Code {
//	return jen.Index().Add(s.Underlying.PbType())
//}
//
//func (s *Slice) MirType() jen.Code {
//	return jen.Index().Add(s.Underlying.MirType())
//}
//
//func (s *Slice) FromPbType(code jen.Code) jen.Code {
//	return jen.Id(s)
//}
