package protoreflectutil

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/filecoin-project/mir/pkg/util/reflectutil"
)

func IsMessageType(goStructPtr reflect.Type) bool {
	return goStructPtr.Kind() == reflect.Pointer &&
		goStructPtr.Implements(reflectutil.TypeOf[protoreflect.ProtoMessage]())
}

func DescriptorForType(goStructPtr reflect.Type) (protoreflect.MessageDescriptor, bool) {
	nilMsg, ok := reflect.Zero(goStructPtr).Interface().(protoreflect.ProtoMessage)
	if !ok {
		return nil, false
	}
	return nilMsg.ProtoReflect().Descriptor(), true
}
