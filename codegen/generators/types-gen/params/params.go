package params

import (
	"strconv"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/generators/types-gen/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// FunctionParam represents a parameter of a function.
type FunctionParam struct {
	ParamName string
	Type      types.Type
}

// MirCode returns the code for the parameter that can be used in a function declaration using Mir-generated
func (p FunctionParam) MirCode() jen.Code {
	return jen.Id(p.ParamName).Add(p.Type.MirType())
}

// ConstructorParam represents a parameter in a constructor of a message.
type ConstructorParam struct {
	FunctionParam
	Field *types.Field
}

// paramName is used in uniqueName.
func (p FunctionParam) paramName() string {
	return p.ParamName
}

// uniqueName returns modifies originalName in such a way that it is different from all names in the list l.
func uniqueName[T interface{ paramName() string }](l []T, originalName string) string {
	suffixNumber := int64(-1)
	nameWithSuffix := originalName
	for _, param := range l {
		if nameWithSuffix == param.paramName() {
			suffixNumber += 1
			nameWithSuffix = originalName + strconv.FormatInt(suffixNumber, 10)
		}
	}

	return nameWithSuffix
}

// FunctionParamList represents a list of parameters of a function.
type FunctionParamList struct {
	Slice []FunctionParam
}

// FunctionParamListOf returns a FunctionParamList containing the given parameters.
func FunctionParamListOf(params ...FunctionParam) FunctionParamList {
	return FunctionParamList{params}
}

// Append returns a new FunctionParamList with an item appended to it.
func (l FunctionParamList) Append(name string, tp types.Type) FunctionParamList {
	param := FunctionParam{
		ParamName: uniqueName(l.Slice, name),
		Type:      tp,
	}

	return FunctionParamList{append(l.Slice, param)}
}

// AppendAll adds all parameters to the list, making sure that all parameters have unique names.
func (l FunctionParamList) AppendAll(other FunctionParamList) FunctionParamList {
	res := l
	for _, param := range other.Slice {
		res = res.Append(param.ParamName, param.Type)
	}

	return res
}

// MirCode returns the slice of function parameters as a jen.Code list to be used in a function declaration.
func (l FunctionParamList) MirCode() []jen.Code {
	return sliceutil.Transform(l.Slice, func(_ int, p FunctionParam) jen.Code { return p.MirCode() })
}

// IDs returns the slice of ParamName of the parameters as jen.Code IDs.
func (l FunctionParamList) IDs() []jen.Code {
	return sliceutil.Transform(l.Slice, func(_ int, p FunctionParam) jen.Code { return jen.Id(p.ParamName) })
}

// ConstructorParamList represents a list of parameters of a constructor function of a message.
type ConstructorParamList struct {
	Slice []ConstructorParam
}

// Append adds a parameter to the list, making sure that all parameters have unique names.
func (l ConstructorParamList) Append(name string, field *types.Field) ConstructorParamList {
	param := ConstructorParam{
		FunctionParam: FunctionParam{
			ParamName: uniqueName(l.Slice, name),
			Type:      field.Type,
		},
		Field: field,
	}

	return ConstructorParamList{append(l.Slice, param)}
}

// AppendAll adds all parameters to the list, making sure that all parameters have unique names.
func (l ConstructorParamList) AppendAll(other ConstructorParamList) ConstructorParamList {
	res := l
	for _, param := range other.Slice {
		res = res.Append(param.ParamName, param.Field)
	}

	return res
}

// MirCode returns the slice of function parameters as a jen.Code list to be used in a function declaration.
func (l ConstructorParamList) MirCode() []jen.Code {
	return sliceutil.Transform(l.Slice, func(_ int, p ConstructorParam) jen.Code { return p.MirCode() })
}

// IDs returns the slice of ParamName of the parameters as jen.Code IDs.
func (l ConstructorParamList) IDs() []jen.Code {
	return sliceutil.Transform(l.Slice, func(_ int, p ConstructorParam) jen.Code { return jen.Id(p.ParamName) })
}

// FunctionParamList transforms ConstructorParamList to FunctionParamList.
func (l ConstructorParamList) FunctionParamList() FunctionParamList {
	return FunctionParamList{
		sliceutil.Transform(l.Slice, func(_ int, p ConstructorParam) FunctionParam {
			return p.FunctionParam
		}),
	}
}
