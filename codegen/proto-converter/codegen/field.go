package codegen

import (
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/astutil"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
	"github.com/filecoin-project/mir/pkg/pb/mir"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// Field represents a field in a protobuf message.
type Field struct {
	// Name is the name of the field.
	Name string
	// PbType is the type of the field in the protobuf-generated struct.
	PbType jen.Code
	// MirType is the type of the field in the Mir-generated struct.
	MirType jen.Code
	// FromPb is a function that converts a protobuf value to a Mir value.
	FromPb jen.Code
	// ToPb is a function that converts a Mir value to a protobuf value.
	ToPb jen.Code
}

// LowercaseName returns the lowercase name of the field.
func (f *Field) LowercaseName() string {
	return astutil.ToUnexported(f.Name)
}

func (f *Field) FuncParamPbType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.PbType)
}

func (f *Field) FuncParamMirType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.MirType)
}

func (f *Field) StructParamPbType() jen.Code {
	return jen.Id(f.Name).Add(f.PbType)
}

func (f *Field) StructParamMirType() jen.Code {
	return jen.Id(f.Name).Add(f.MirType)
}

// GetField extracts the information about the field necessary for code generation.
func GetField(goField reflect.StructField, protoField protoreflect.FieldDescriptor) (*Field, error) {
	// If the field is a Mir-annotated message, use the Mir type.
	if protoField.Kind() == protoreflect.MessageKind {
		msg, err := MessageFromPbGoType(goField.Type)
		if err != nil {
			return nil, err
		}
		if msg.ShouldGenerateMirType() {
			return &Field{
				Name:    goField.Name,
				PbType:  msg.PbTypePtr(),
				MirType: msg.MirTypePtr(),
				ToPb:    msg.ToPbFunc(),
				FromPb:  msg.FromPbFunc(),
			}, nil
		}
	}

	// Otherwise, check if the field has (mir.type) option specified.
	protoFieldOptions := protoField.Options().(*descriptorpb.FieldOptions)
	mirTypeOption := proto.GetExtension(protoFieldOptions, mir.E_Type).(string)

	// If the option is specified, use it. Otherwise, use pb type.
	pbType := jenutil.QualFromType(goField.Type)
	var mirType jen.Code
	if mirTypeOption != "" {
		sepIdx := strings.LastIndex(mirTypeOption, ".")
		mirType = jen.Qual(mirTypeOption[:sepIdx], mirTypeOption[sepIdx+1:])
	} else {
		mirType = pbType
	}

	return &Field{
		Name:    goField.Name,
		PbType:  pbType,
		MirType: mirType,
		ToPb:    jen.Parens(pbType),
		FromPb:  jen.Parens(mirType),
	}, nil
}

type Fields []*Field

func (fs Fields) FuncParamsPbTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.FuncParamPbType() })
}

func (fs Fields) FuncParamsMirTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.FuncParamMirType() })
}

func (fs Fields) StructParamsPbTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.StructParamPbType() })
}

func (fs Fields) StructParamsMirTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.StructParamMirType() })
}
