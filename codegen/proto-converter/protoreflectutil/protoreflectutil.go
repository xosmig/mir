package protoreflectutil

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/filecoin-project/mir/pkg/util/reflectutil"
)

func IsMessageStruct(tp reflect.Type) bool {
	return tp.Kind() == reflect.Struct &&
		reflect.PointerTo(tp).Implements(reflectutil.TypeOf[protoreflect.ProtoMessage]())
}

func DescriptorForType(tp reflect.Type) (protoreflect.MessageDescriptor, bool) {
	msg, ok := reflect.Zero(tp).Addr().Interface().(protoreflect.ProtoMessage)
	if !ok {
		return nil, false
	}
	return msg.ProtoReflect().Descriptor(), true
}
