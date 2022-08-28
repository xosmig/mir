package model

import (
	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/astutil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// Field represents a field in a protobuf message.
type Field struct {
	// Name is the name of the field .
	Name string

	// Type contains type-related information about the field.
	Type Type
}

// LowercaseName returns the lowercase name of the field.
func (f *Field) LowercaseName() string {
	return astutil.ToUnexported(f.Name)
}

func (f *Field) FuncParamPbType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.Type.PbType())
}

func (f *Field) FuncParamMirType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.Type.MirType())
}

type Fields []*Field

func (fs Fields) FuncParamsPbTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.FuncParamPbType() })
}

func (fs Fields) FuncParamsMirTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.FuncParamMirType() })
}

func (fs Fields) FuncParamsIDs() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return jen.Id(f.Name) })
}
