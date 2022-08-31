package events

import (
	"strconv"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/model/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// ConstructorParam represents a parameter in a constructor of a message.
type ConstructorParam struct {
	ParamName string
	Field     *types.Field
}

// MirCode returns the code for the parameter that can be used in a function declaration using Mir-generated types.
func (p ConstructorParam) MirCode() jen.Code {
	return jen.Id(p.ParamName).Add(p.Field.Type.MirType())
}

// ConstructorParamList represents a list of parameters of a constructor function of a message.
type ConstructorParamList struct {
	Slice []ConstructorParam
}

// Append adds a parameter to the list, making sure that all parameters have unique names.
func (l ConstructorParamList) Append(name string, field *types.Field) ConstructorParamList {
	suffixNumber := int64(-1)
	nameWithSuffix := name
	for _, param := range l.Slice {
		if nameWithSuffix == param.ParamName {
			suffixNumber += 1
			nameWithSuffix = name + strconv.FormatInt(suffixNumber, 10)
		}
	}

	return ConstructorParamList{append(l.Slice, ConstructorParam{nameWithSuffix, field})}
}

// MirCode returns the slice of function parameters as a jen.Code list to be used in a function declaration.
func (l ConstructorParamList) MirCode() []jen.Code {
	return sliceutil.Transform(l.Slice, func(_ int, p ConstructorParam) jen.Code { return p.MirCode() })
}

// IDs returns the slice of ParamName of the parameters as jen.Code IDs.
func (l ConstructorParamList) IDs() []jen.Code {
	return sliceutil.Transform(l.Slice, func(_ int, p ConstructorParam) jen.Code { return jen.Id(p.ParamName) })
}
