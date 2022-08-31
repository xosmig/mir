package jenutil

import (
	"strconv"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// FuncParam represents a function parameter.
type FuncParam struct {
	LowercaseName string
	Type          jen.Code
}

// Code returns the code for the parameter that can be used in a function declaration.
func (p FuncParam) Code() jen.Code {
	return jen.Id(p.LowercaseName).Add(p.Type)
}

// FuncParamList represents a list of function parameters.
type FuncParamList struct {
	slice []FuncParam
}

func NewFuncParamList() *FuncParamList {
	return &FuncParamList{}
}

// Append adds a parameter to the list, making sure that all parameters have unique names.
func (l *FuncParamList) Append(name string, tp jen.Code) {
	suffixNumber := int64(-1)
	nameWithSuffix := name
	for _, param := range l.slice {
		if nameWithSuffix == param.LowercaseName {
			suffixNumber += 1
			nameWithSuffix = name + strconv.FormatInt(suffixNumber, 10)
		}
	}

	l.slice = append(l.slice, FuncParam{nameWithSuffix, tp})
}

// Slice returns the slice of function parameters.
func (l *FuncParamList) Slice() []FuncParam {
	return l.slice
}

// Code returns the slice of function parameters as a jen.Code list to be used in a function declaration.
func (l *FuncParamList) Code() []jen.Code {
	return sliceutil.Transform(l.slice, func(_ int, p FuncParam) jen.Code { return p.Code() })
}
