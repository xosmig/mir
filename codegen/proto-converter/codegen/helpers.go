package codegen

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/pkg/pb/mir"
)

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
