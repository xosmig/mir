package model

import (
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
)

type Oneof struct {
	Name   string
	Parent *Message
	// To simplify parsing, Oneof doesn't store a list of its options.
}

func (t *Oneof) Same() bool {
	return t.Parent.Same()
}

func (t *Oneof) PbExportedInterfaceName() string {
	return fmt.Sprintf("%v_%v", t.Parent.Name(), t.Name)
}

func (t *Oneof) PbNativeInterfaceName() string {
	return fmt.Sprintf("is%v_%v", t.Parent.Name(), t.Name)
}

func (t *Oneof) MirInterfaceName() string {
	return fmt.Sprintf("%v_%v", t.Parent.Name(), t.Name)
}

func (t *Oneof) PbMethodName() string {
	return "is" + t.PbExportedInterfaceName()
}

func (t *Oneof) MirMethodName() string {
	return t.PbMethodName()
}

func (t *Oneof) PbType() jen.Code {
	return jen.Qual(t.Parent.PbPkgPath(), t.PbExportedInterfaceName())
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

func (opt *OneofOption) PbWrapperType() jen.Code {
	return jen.Op("*").Add(jenutil.QualFromType(opt.PbWrapperReflect.Elem()))
}

func (opt *OneofOption) NewPbWrapperType() jen.Code {
	return jen.Op("&").Add(jenutil.QualFromType(opt.PbWrapperReflect.Elem()))
}

func (opt *OneofOption) MirWrapperStructType() jen.Code {
	return jen.Qual(StructsPackagePath(opt.PbWrapperReflect.Elem().PkgPath()), opt.WrapperName)
}

func (opt *OneofOption) MirWrapperType() jen.Code {
	return jen.Op("*").Add(opt.MirWrapperStructType())
}

func (opt *OneofOption) NewMirWrapperType() jen.Code {
	return jen.Op("&").Add(opt.MirWrapperStructType())
}
