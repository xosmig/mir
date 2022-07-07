// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: mir/plugin.proto

package mir

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_mir_plugin_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.OneofOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50000,
		Name:          "mir.oneof",
		Tag:           "varint,50000,opt,name=oneof",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.OneofOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50001,
		Name:          "mir.event_type",
		Tag:           "varint,50001,opt,name=event_type",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50000,
		Name:          "mir.struct",
		Tag:           "varint,50000,opt,name=struct",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50001,
		Name:          "mir.event",
		Tag:           "varint,50001,opt,name=event",
		Filename:      "mir/plugin.proto",
	},
}

// Extension fields to descriptorpb.OneofOptions.
var (
	// optional bool oneof = 50000;
	E_Oneof = &file_mir_plugin_proto_extTypes[0]
	// optional bool event_type = 50001;
	E_EventType = &file_mir_plugin_proto_extTypes[1]
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional bool struct = 50000;
	E_Struct = &file_mir_plugin_proto_extTypes[2]
	// optional bool event = 50001;
	E_Event = &file_mir_plugin_proto_extTypes[3]
)

var File_mir_plugin_proto protoreflect.FileDescriptor

var file_mir_plugin_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x6d, 0x69, 0x72, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x35, 0x0a, 0x05, 0x6f, 0x6e, 0x65,
	0x6f, 0x66, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x6f, 0x6e, 0x65, 0x6f, 0x66,
	0x3a, 0x3e, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x3a, 0x39, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x3a, 0x37, 0x0a, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x6d,
	0x69, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_mir_plugin_proto_goTypes = []interface{}{
	(*descriptorpb.OneofOptions)(nil),   // 0: google.protobuf.OneofOptions
	(*descriptorpb.MessageOptions)(nil), // 1: google.protobuf.MessageOptions
}
var file_mir_plugin_proto_depIdxs = []int32{
	0, // 0: mir.oneof:extendee -> google.protobuf.OneofOptions
	0, // 1: mir.event_type:extendee -> google.protobuf.OneofOptions
	1, // 2: mir.struct:extendee -> google.protobuf.MessageOptions
	1, // 3: mir.event:extendee -> google.protobuf.MessageOptions
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	0, // [0:4] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mir_plugin_proto_init() }
func file_mir_plugin_proto_init() {
	if File_mir_plugin_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mir_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_mir_plugin_proto_goTypes,
		DependencyIndexes: file_mir_plugin_proto_depIdxs,
		ExtensionInfos:    file_mir_plugin_proto_extTypes,
	}.Build()
	File_mir_plugin_proto = out.File
	file_mir_plugin_proto_rawDesc = nil
	file_mir_plugin_proto_goTypes = nil
	file_mir_plugin_proto_depIdxs = nil
}
