package main

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/pb/mir"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
)

type GoType struct {
	UserType        string
	PbType          string
	NeedsConversion bool
}

func GetGoType(g *protogen.GeneratedFile, field *protogen.Field) (GoType, error) {
	pbType, err := GetPbType(field)
	if err != nil {
		return GoType{}, err
	}

	fieldOptions := field.Desc.Options().(*descriptorpb.FieldOptions)
	goTypeExt := proto.GetExtension(fieldOptions, mir.E_TypeWrapper).(string)
	if goTypeExt != "" {
		ident, err := GoIdentFromString(goTypeExt)
		if err != nil {
			return GoType{}, err
		}
		return GoType{
			UserType:        g.QualifiedGoIdent(ident),
			PbType:          pbType,
			NeedsConversion: true,
		}, nil
	}

	return GoType{
		UserType:        pbType,
		PbType:          pbType,
		NeedsConversion: false,
	}, nil
}

//func GetPbType(field *protogen.Field) (string, error) {
//	if field.Desc.Message() != nil {
//		return string(field.Desc.Message().FullName()), nil
//	}
//	return field.Desc.Kind().String(), nil
//}

func GetPbType(field *protogen.Field) (string, error) {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return "bool", nil
	case protoreflect.EnumKind:
		return "", fmt.Errorf("enum fields in events are not yet supported by the plugin")
	case protoreflect.Int32Kind, protoreflect.Sint32Kind:
		return "int32", nil
	case protoreflect.Uint32Kind:
		return "uint32", nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind:
		return "int64", nil
	case protoreflect.Uint64Kind:
		return "uint64", nil
	case protoreflect.Sfixed32Kind:
		// TODO
		return "", fmt.Errorf("TODO")
	case protoreflect.Fixed32Kind:
		// TODO
		return "", fmt.Errorf("TODO")
	case protoreflect.FloatKind:
		return "float32", nil
	case protoreflect.Sfixed64Kind:
		// TODO
		return "", fmt.Errorf("TODO")
	case protoreflect.Fixed64Kind:
		// TODO
		return "", fmt.Errorf("TODO")
	case protoreflect.DoubleKind:
		return "float64", nil
	case protoreflect.StringKind:
		return "string", nil
	case protoreflect.BytesKind:
		return "[]byte", nil
	case protoreflect.MessageKind:
		// TODO
		return "", fmt.Errorf("TODO")
	case protoreflect.GroupKind:
		// TODO
		return "", fmt.Errorf("TODO")
	default:
		return "", fmt.Errorf("unknown field kind %v", field.Desc.Kind())
	}
}

func GoIdentFromString(s string) (protogen.GoIdent, error) {
	delimeterIdx := strings.LastIndex(s, ".")
	if delimeterIdx == -1 || delimeterIdx == 0 || delimeterIdx == len(s)-1 {
		return protogen.GoIdent{},
			fmt.Errorf("invalid type identified: %v. Expected format: full/package/name.TypeName", s)
	}

	return protogen.GoIdent{
		GoName:       s[delimeterIdx+1:],
		GoImportPath: protogen.GoImportPath(s[:delimeterIdx]),
	}, nil
}
