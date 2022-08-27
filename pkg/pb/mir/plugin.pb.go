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
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51200,
		Name:          "mir.event",
		Tag:           "varint,51200,opt,name=event",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51201,
		Name:          "mir.message",
		Tag:           "varint,51201,opt,name=message",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51202,
		Name:          "mir.struct",
		Tag:           "varint,51202,opt,name=struct",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51203,
		Name:          "mir.parent_event",
		Tag:           "bytes,51203,opt,name=parent_event",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51204,
		Name:          "mir.event_root",
		Tag:           "varint,51204,opt,name=event_root",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51205,
		Name:          "mir.parent_oneof",
		Tag:           "bytes,51205,opt,name=parent_oneof",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51206,
		Name:          "mir.oneof_wrapper",
		Tag:           "bytes,51206,opt,name=oneof_wrapper",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.OneofOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         51200,
		Name:          "mir.event_type",
		Tag:           "varint,51200,opt,name=event_type",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         51200,
		Name:          "mir.type",
		Tag:           "bytes,51200,opt,name=type",
		Filename:      "mir/plugin.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         51200,
		Name:          "mir.imports",
		Tag:           "bytes,51200,rep,name=imports",
		Filename:      "mir/plugin.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional bool event = 51200;
	E_Event = &file_mir_plugin_proto_extTypes[0]
	// optional bool message = 51201;
	E_Message = &file_mir_plugin_proto_extTypes[1]
	// optional bool struct = 51202;
	E_Struct = &file_mir_plugin_proto_extTypes[2]
	// optional string parent_event = 51203;
	E_ParentEvent = &file_mir_plugin_proto_extTypes[3]
	// optional bool event_root = 51204;
	E_EventRoot = &file_mir_plugin_proto_extTypes[4]
	// optional string parent_oneof = 51205;
	E_ParentOneof = &file_mir_plugin_proto_extTypes[5]
	// optional string oneof_wrapper = 51206;
	E_OneofWrapper = &file_mir_plugin_proto_extTypes[6]
)

// Extension fields to descriptorpb.OneofOptions.
var (
	// optional bool event_type = 51200;
	E_EventType = &file_mir_plugin_proto_extTypes[7]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional string type = 51200;
	E_Type = &file_mir_plugin_proto_extTypes[8]
)

// Extension fields to descriptorpb.FileOptions.
var (
	// repeated string imports = 51200;
	E_Imports = &file_mir_plugin_proto_extTypes[9]
)

var File_mir_plugin_proto protoreflect.FileDescriptor

var file_mir_plugin_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x6d, 0x69, 0x72, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x37, 0x0a, 0x05, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x80, 0x90, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x3a, 0x3b, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x81,
	0x90, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x3a,
	0x39, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x82, 0x90, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x3a, 0x44, 0x0a, 0x0c, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x83, 0x90, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x3a, 0x40, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x6f, 0x6f, 0x74, 0x12, 0x1f,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x84, 0x90, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x6f,
	0x6f, 0x74, 0x3a, 0x44, 0x0a, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6f, 0x6e, 0x65,
	0x6f, 0x66, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x85, 0x90, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x3a, 0x46, 0x0a, 0x0d, 0x6f, 0x6e, 0x65, 0x6f,
	0x66, 0x5f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x86, 0x90, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x3a, 0x3e, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x80, 0x90,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x3a, 0x33, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x80, 0x90, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x38, 0x0a, 0x07, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73,
	0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x80,
	0x90, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x42,
	0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d,
	0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x6d, 0x69, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_mir_plugin_proto_goTypes = []interface{}{
	(*descriptorpb.MessageOptions)(nil), // 0: google.protobuf.MessageOptions
	(*descriptorpb.OneofOptions)(nil),   // 1: google.protobuf.OneofOptions
	(*descriptorpb.FieldOptions)(nil),   // 2: google.protobuf.FieldOptions
	(*descriptorpb.FileOptions)(nil),    // 3: google.protobuf.FileOptions
}
var file_mir_plugin_proto_depIdxs = []int32{
	0,  // 0: mir.event:extendee -> google.protobuf.MessageOptions
	0,  // 1: mir.message:extendee -> google.protobuf.MessageOptions
	0,  // 2: mir.struct:extendee -> google.protobuf.MessageOptions
	0,  // 3: mir.parent_event:extendee -> google.protobuf.MessageOptions
	0,  // 4: mir.event_root:extendee -> google.protobuf.MessageOptions
	0,  // 5: mir.parent_oneof:extendee -> google.protobuf.MessageOptions
	0,  // 6: mir.oneof_wrapper:extendee -> google.protobuf.MessageOptions
	1,  // 7: mir.event_type:extendee -> google.protobuf.OneofOptions
	2,  // 8: mir.type:extendee -> google.protobuf.FieldOptions
	3,  // 9: mir.imports:extendee -> google.protobuf.FileOptions
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	0,  // [0:10] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
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
			NumExtensions: 10,
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
