package model

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/astutil"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/protoreflectutil"
	"github.com/filecoin-project/mir/pkg/pb/mir"
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

// GetField extracts the information about the field necessary for code generation.
func GetField(goField reflect.StructField, protoField protoreflect.FieldDescriptor) (*Field, error) {
	tp, err := getFieldType(goField.Type, protoField)
	if err != nil {
		return nil, err
	}

	return &Field{
		Name: goField.Name,
		Type: tp,
	}, nil
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

func getFieldType(goType reflect.Type, protoField protoreflect.FieldDescriptor) (Type, error) {
	// TODO: Since maps are not currently used, I didn't bother supporting them yet.
	if goType.Kind() == reflect.Map {
		return nil, fmt.Errorf("map fields are not supported yet")
	}

	// Check if the field is repeated.
	if goType.Kind() == reflect.Slice {
		underlying, err := getFieldType(goType.Elem(), protoField)
		if err != nil {
			return nil, err
		}
		return Slice{underlying}, nil
	}

	// Check if the field has (mir.type) option specified.
	protoFieldOptions := protoField.Options().(*descriptorpb.FieldOptions)
	mirTypeOption := proto.GetExtension(protoFieldOptions, mir.E_Type).(string)
	if mirTypeOption != "" {
		sepIdx := strings.LastIndex(mirTypeOption, ".")
		return Castable{
			pbType:  jenutil.QualFromType(goType),
			mirType: jen.Qual(mirTypeOption[:sepIdx], mirTypeOption[sepIdx+1:]),
		}, nil
	}

	// Check if the field is a message.
	if protoreflectutil.IsMessageType(goType) {
		msg, err := ParseMessage(goType)
		if err != nil {
			return nil, err
		}

		return msg, nil
	}

	return Same{jenutil.QualFromType(goType)}, nil
}
