package codegen

import (
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

type Field struct {
	lowercaseName string

	Name    string
	PbType  jen.Code
	MirType jen.Code
}

func (f Field) LowercaseName() string {
	if f.lowercaseName == "" && f.Name != "" {
		f.lowercaseName = strings.ToLower(f.Name[:1]) + f.Name[1:]
	}
	return f.lowercaseName
}

func (f Field) FuncParamPbType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.PbType)
}

func (f Field) FuncParamMirType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.MirType)
}

type Fields []Field

func (fs Fields) FuncParamsPbTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, t Field) jen.Code {
		return t.FuncParamPbType()
	})
}

func (fs Fields) FuncParamsMirTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, t Field) jen.Code {
		return t.FuncParamMirType()
	})
}
