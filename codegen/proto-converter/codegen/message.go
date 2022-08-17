package codegen

import (
	"go/ast"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/codegen/proto-converter/jenutil"
	"github.com/filecoin-project/mir/codegen/proto-converter/protoreflectutil"
	"github.com/filecoin-project/mir/pkg/pb/mir"
)

type Message struct {
	PbGoStruct reflect.Type
	ProtoDesc  protoreflect.MessageDescriptor
	Fields     Fields
}

func MessageFromPbStructType(pbGoStruct reflect.Type) (*Message, bool, error) {
	msgDesc, ok := protoreflectutil.DescriptorForType(pbGoStruct)
	if !ok {
		// This struct is not a ProtoMessage.
		return nil, false, nil
	}

	var fields Fields

	for i := 0; i < pbGoStruct.NumField(); i++ {
		// Get go representation of the field.
		goField := pbGoStruct.Field(i)
		if !ast.IsExported(goField.Name) {
			// Skip unexported fields.
			continue
		}

		if oneofTag, ok := goField.Tag.Lookup("protobuf_oneof"); ok {
			_ = oneofTag
			// TODO: oneofs are skipped for now.
			continue
		}

		// Get protobuf representation of the field.
		protoName, err := getProtoNameOfField(goField)
		if err != nil {
			return nil, false, err
		}
		protoField := msgDesc.Fields().ByName(protoreflect.Name(protoName))

		// Get the mir type for the field.
		mirType, err := getMirType(goField, protoField)
		if err != nil {
			return nil, false, err
		}

		fields = append(fields, Field{
			Name:    goField.Name,
			PbType:  jenutil.QualFromType(goField.Type),
			MirType: mirType,
		})
	}

	return &Message{
		PbGoStruct: pbGoStruct,
		ProtoDesc:  msgDesc,
		Fields:     fields,
	}, true, nil
}

func (m *Message) IsMirEvent() bool {
	return proto.GetExtension(m.ProtoDesc.Options().(*descriptorpb.MessageOptions), mir.E_Event).(bool)
}

func (m *Message) IsMirMessage() bool {
	return proto.GetExtension(m.ProtoDesc.Options().(*descriptorpb.MessageOptions), mir.E_Message).(bool)
}

func (m *Message) IsMirStruct() bool {
	return proto.GetExtension(m.ProtoDesc.Options().(*descriptorpb.MessageOptions), mir.E_Struct).(bool)
}
