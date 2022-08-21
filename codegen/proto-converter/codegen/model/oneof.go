package model

import (
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
)

type Oneof struct {
	Name    string
	Parent  *Message
	Options []*OneofOption
}

func (t *Oneof) Same() bool {
	return t.Parent.Same()
}

func (t *Oneof) PbInterfaceName() string {
	return fmt.Sprintf("%v_%v", t.Parent.Name(), t.Name)
}

func (t *Oneof) MirInterfaceName() string {
	return fmt.Sprintf("%v_%v", t.Parent.Name(), t.Name)
}

func (t *Oneof) PbMethodName() string {
	return "is" + t.PbInterfaceName()
}

func (t *Oneof) MirMethodName() string {
	return t.PbMethodName()
}

func (t *Oneof) PbType() jen.Code {
	return jen.Qual(t.Parent.PbPkgPath(), t.PbInterfaceName())
}

func (t *Oneof) MirType() jen.Code {
	return jen.Qual(t.Parent.MirPkgPath(), t.MirInterfaceName())
}

func (t *Oneof) ToMir(code jen.Code) jen.Code {
	return jen.Qual(t.Parent.MirPkgPath(), t.MirInterfaceName()+"FromPb").Call(code)
}

func (t *Oneof) ToPb(code jen.Code) jen.Code {
	return jen.Add(code).Dot("Pb").Call()
}

type OneofOption struct {
	PbWrapperReflect reflect.Type
	WrapperName      string
	Field            *Field
}

func (opt *OneofOption) NewPbWrapperType() jen.Code {
	return jen.Op("&").Add(jenutil.QualFromType(opt.PbWrapperReflect.Elem()))
}
