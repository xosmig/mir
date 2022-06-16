// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: cbpb/cbpb.proto

package cbpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CBMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//	*CBMessage_StartMessage
	//	*CBMessage_EchoMessage
	Type isCBMessage_Type `protobuf_oneof:"type"`
}

func (x *CBMessage) Reset() {
	*x = CBMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cbpb_cbpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CBMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CBMessage) ProtoMessage() {}

func (x *CBMessage) ProtoReflect() protoreflect.Message {
	mi := &file_cbpb_cbpb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CBMessage.ProtoReflect.Descriptor instead.
func (*CBMessage) Descriptor() ([]byte, []int) {
	return file_cbpb_cbpb_proto_rawDescGZIP(), []int{0}
}

func (m *CBMessage) GetType() isCBMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *CBMessage) GetStartMessage() *StartMessage {
	if x, ok := x.GetType().(*CBMessage_StartMessage); ok {
		return x.StartMessage
	}
	return nil
}

func (x *CBMessage) GetEchoMessage() *EchoMessage {
	if x, ok := x.GetType().(*CBMessage_EchoMessage); ok {
		return x.EchoMessage
	}
	return nil
}

type isCBMessage_Type interface {
	isCBMessage_Type()
}

type CBMessage_StartMessage struct {
	StartMessage *StartMessage `protobuf:"bytes,1,opt,name=startMessage,proto3,oneof"`
}

type CBMessage_EchoMessage struct {
	EchoMessage *EchoMessage `protobuf:"bytes,2,opt,name=echoMessage,proto3,oneof"`
}

func (*CBMessage_StartMessage) isCBMessage_Type() {}

func (*CBMessage_EchoMessage) isCBMessage_Type() {}

type StartMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *StartMessage) Reset() {
	*x = StartMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cbpb_cbpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartMessage) ProtoMessage() {}

func (x *StartMessage) ProtoReflect() protoreflect.Message {
	mi := &file_cbpb_cbpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartMessage.ProtoReflect.Descriptor instead.
func (*StartMessage) Descriptor() ([]byte, []int) {
	return file_cbpb_cbpb_proto_rawDescGZIP(), []int{1}
}

func (x *StartMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type EchoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *EchoMessage) Reset() {
	*x = EchoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cbpb_cbpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoMessage) ProtoMessage() {}

func (x *EchoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_cbpb_cbpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoMessage.ProtoReflect.Descriptor instead.
func (*EchoMessage) Descriptor() ([]byte, []int) {
	return file_cbpb_cbpb_proto_rawDescGZIP(), []int{2}
}

func (x *EchoMessage) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_cbpb_cbpb_proto protoreflect.FileDescriptor

var file_cbpb_cbpb_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x62, 0x70, 0x62, 0x2f, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x63, 0x62, 0x70, 0x62, 0x22, 0x84, 0x01, 0x0a, 0x09, 0x43, 0x42, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x38, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x62,
	0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48,
	0x00, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x35, 0x0a, 0x0b, 0x65, 0x63, 0x68, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x62, 0x70, 0x62, 0x2e, 0x45, 0x63, 0x68, 0x6f,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x65, 0x63, 0x68, 0x6f, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x22,
	0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x21, 0x0a, 0x0b, 0x45, 0x63, 0x68, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f,
	0x63, 0x62, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cbpb_cbpb_proto_rawDescOnce sync.Once
	file_cbpb_cbpb_proto_rawDescData = file_cbpb_cbpb_proto_rawDesc
)

func file_cbpb_cbpb_proto_rawDescGZIP() []byte {
	file_cbpb_cbpb_proto_rawDescOnce.Do(func() {
		file_cbpb_cbpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_cbpb_cbpb_proto_rawDescData)
	})
	return file_cbpb_cbpb_proto_rawDescData
}

var file_cbpb_cbpb_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_cbpb_cbpb_proto_goTypes = []interface{}{
	(*CBMessage)(nil),    // 0: cbpb.CBMessage
	(*StartMessage)(nil), // 1: cbpb.StartMessage
	(*EchoMessage)(nil),  // 2: cbpb.EchoMessage
}
var file_cbpb_cbpb_proto_depIdxs = []int32{
	1, // 0: cbpb.CBMessage.startMessage:type_name -> cbpb.StartMessage
	2, // 1: cbpb.CBMessage.echoMessage:type_name -> cbpb.EchoMessage
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_cbpb_cbpb_proto_init() }
func file_cbpb_cbpb_proto_init() {
	if File_cbpb_cbpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cbpb_cbpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CBMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cbpb_cbpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cbpb_cbpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_cbpb_cbpb_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*CBMessage_StartMessage)(nil),
		(*CBMessage_EchoMessage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cbpb_cbpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cbpb_cbpb_proto_goTypes,
		DependencyIndexes: file_cbpb_cbpb_proto_depIdxs,
		MessageInfos:      file_cbpb_cbpb_proto_msgTypes,
	}.Build()
	File_cbpb_cbpb_proto = out.File
	file_cbpb_cbpb_proto_rawDesc = nil
	file_cbpb_cbpb_proto_goTypes = nil
	file_cbpb_cbpb_proto_depIdxs = nil
}
