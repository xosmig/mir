package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GoType(field *protogen.Field) (string, error) {
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
