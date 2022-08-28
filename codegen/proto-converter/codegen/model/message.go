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
	"github.com/filecoin-project/mir/pkg/pb/mir"
)

func StructsPackageName(pbPackagePath string) string {
	return pbPackagePath[strings.LastIndex(pbPackagePath, "/")+1:] + "structs"
}

func StructsPackagePath(pbPackagePath string) string {
	return fmt.Sprintf("%v/%v", pbPackagePath, StructsPackageName(pbPackagePath))
}

// Message contains the information needed to generate code for a protobuf message.
type Message struct {
	shouldGenerateMirType bool
	mirPkgPath            string

	fields               Fields
	pbStructType         jen.Code
	mirStructType        jen.Code
	protoDesc            protoreflect.MessageDescriptor
	pbGoStructPtrReflect reflect.Type
}

func (m *Message) Name() string {
	return m.pbGoStructPtrReflect.Elem().Name()
}

func (m *Message) PbPkgPath() string {
	return m.pbGoStructPtrReflect.Elem().PkgPath()
}

func (m *Message) MirPkgPath() string {
	// Return the cached value if present.
	if m.mirPkgPath != "" {
		return m.mirPkgPath
	}

	m.mirPkgPath = StructsPackagePath(m.PbPkgPath())
	return m.mirPkgPath
}

func (m *Message) Same() bool {
	return !m.ShouldGenerateMirType()
}

func (m *Message) PbType() jen.Code {
	return jen.Op("*").Add(m.pbStructType)
}

func (m *Message) NewPbType() jen.Code {
	return jen.Op("&").Add(m.pbStructType)
}

func (m *Message) NewMirType() jen.Code {
	return jen.Op("&").Add(m.mirStructType)
}

func (m *Message) MirType() jen.Code {
	return jen.Op("*").Add(m.mirStructType)
}

func (m *Message) ToMir(code jen.Code) jen.Code {
	if m.Same() {
		return code
	}
	return jen.Qual(m.MirPkgPath(), m.Name()+"FromPb").Call(code)
}

func (m *Message) ToPb(code jen.Code) jen.Code {
	if m.Same() {
		return code
	}
	return jen.Add(code).Dot("Pb").Call()
}

func (m *Message) ConstructorName() string {
	return "New" + m.Name()
}

func (m *Message) Constructor() *jen.Statement {
	return jen.Qual(m.MirPkgPath(), m.ConstructorName())
}

// LowercaseName returns the name of the message in lowercase.
func (m *Message) LowercaseName() string {
	return astutil.ToUnexported(m.Name())
}

//func (m *Message) FuncParamPbType() jen.Code {
//	return jen.Id(m.LowercaseName()).Add(m.PbType())
//}
//
//func (m *Message) FuncParamMirType() jen.Code {
//	return jen.Id(m.LowercaseName()).Add(m.MirType())
//}
//
//func (m *Message) StructParamPbType() jen.Code {
//	return jen.Id(m.Name()).Add(m.PbType())
//}
//
//func (m *Message) StructParamMirType() jen.Code {
//	return jen.Id(m.Name()).Add(m.MirType())
//}

func (m *Message) IsMirEvent() bool {
	return IsMirEvent(m.protoDesc)
}

func (m *Message) IsMirMessage() bool {
	return IsMirMessage(m.protoDesc)
}

func (m *Message) IsMirStruct() bool {
	return IsMirStruct(m.protoDesc)
}

// ShouldGenerateMirType returns true if Mir should generate a struct for the message type.
func (m *Message) ShouldGenerateMirType() bool {
	return m.shouldGenerateMirType
}

//func (m *Message) ParentEvent() (*Message, error) {
//	// Return the cached version if present.
//	if m.parentEvent != nil {
//		return m.parentEvent, nil
//	}
//
//	ext := proto.GetExtension(m.protoDesc.Options().(*descriptorpb.MessageOptions), mir.E_ParentEvent).(string)
//	if ext == "" {
//		return nil, nil
//	}
//
//	sepIdx := strings.LastIndex(ext, ".")
//	parentPackage, parentType := ext[:sepIdx], ext[sepIdx+1:]
//	if parentType == "" {
//		return nil, fmt.Errorf("invalid format for option (mir.parent_event)")
//	}
//
//	// If the parent package is not specified, the same package is used.
//	if parentPackage == "" {
//		parentPackage = m.PbPkgPath()
//	}
//
//	return
//}

func IsMirEvent(protoDesc protoreflect.MessageDescriptor) bool {
	return proto.GetExtension(protoDesc.Options().(*descriptorpb.MessageOptions), mir.E_Event).(bool)
}

func IsMirMessage(protoDesc protoreflect.MessageDescriptor) bool {
	return proto.GetExtension(protoDesc.Options().(*descriptorpb.MessageOptions), mir.E_Message).(bool)
}

func IsMirStruct(protoDesc protoreflect.MessageDescriptor) bool {
	return proto.GetExtension(protoDesc.Options().(*descriptorpb.MessageOptions), mir.E_Struct).(bool)
}

func ShouldGenerateMirType(protoDesc protoreflect.MessageDescriptor) bool {
	return IsMirEvent(protoDesc) || IsMirMessage(protoDesc) || IsMirStruct(protoDesc)
}

func getProtoNameOfField(field reflect.StructField) (protoName string, err error) {
	protobufTag, ok := field.Tag.Lookup("protobuf")
	if !ok {
		return "", fmt.Errorf("field %v has no protobuf tag", field.Name)
	}

	for _, tagPart := range strings.Split(protobufTag, ",") {
		if strings.HasPrefix(tagPart, "name=") {
			return strings.TrimPrefix(tagPart, "name="), nil
		}
	}

	return "", fmt.Errorf("proto name of field %v is not specified in the tag", field.Name)
}

//func GetParentOneof(protoDesc protoreflect.MessageDescriptor) (*Oneof, error) {
//	str := proto.GetExtension(protoDesc.Options().(*descriptorpb.MessageOptions), mir.E_ParentOneof).(string)
//	parts := strings.Split(str, ".")
//	if len(parts) < 3 {
//		return nil, fmt.Errorf("invalid format for parent_oneof option. " +
//			"Expected format: \"full/pkg/name.Message.oneof_field\"")
//	}
//
//	parentMsgName := parts[len(parts)-2]
//	oneofName := parts[len(parts)-1]
//	parentPkg := strings.Join(parts[:len(parts)-2], "/")
//
//
//
//	return &Oneof{
//		Name: oneofName,
//		Parent:
//	}
//
//}
