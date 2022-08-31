package model

import (
	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/astutil"
	"github.com/filecoin-project/mir/pkg/pb/mir"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// Field represents a field in a protobuf message.
// Note: oneofs are not considered as fields in the protobuf data module, but are considered as fields here.
// The reason for that is that oneofs are mapped to fields in the generated Go code.
type Field struct {
	// The name of the field.
	Name string

	// The information about the type of the field.
	Type Type

	// The protobuf descriptor of the field.
	// The descriptor can be either protoreflect.FieldDescriptor or protoreflect.OneofDescriptor.
	ProtoDesc protoreflect.Descriptor
}

// LowercaseName returns the lowercase name of the field.
func (f *Field) LowercaseName() string {
	return astutil.ToUnexported(f.Name)
}

// FuncParamPbType returns the field lowercase name followed by its pb type.
func (f *Field) FuncParamPbType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.Type.PbType())
}

// FuncParamMirType returns the field lowercase name followed by its mir type.
func (f *Field) FuncParamMirType() jen.Code {
	return jen.Id(f.LowercaseName()).Add(f.Type.MirType())
}

// OmitInConstructor returns true iff the field is marked with option [(mir.omit_in_constructor) = true].
func (f *Field) OmitInConstructor() bool {
	fieldDesc, ok := f.ProtoDesc.(protoreflect.FieldDescriptor)
	if !ok {
		// This is a oneof and this option does not exist for oneofs.
		// Return false as the default value.
		return false
	}

	return proto.GetExtension(fieldDesc.Options().(*descriptorpb.FieldOptions), mir.E_OmitInConstructor).(bool)
}

// IsEventTypeOneof returns true iff the field is a oneof and either
// (a) it is named `type` or `Type`; or
// (b) it is marked with `option (mir.event_type) = true`.
//
// The reason why option (a) is present is that it is easy to forget the annotation and get some weird
// unexpected results.
func (f *Field) IsEventTypeOneof() bool {
	oneofDesc, ok := f.ProtoDesc.(protoreflect.OneofDescriptor)
	if !ok {
		return false
	}
	if f.Name == "Type" {
		return true
	}

	return proto.GetExtension(oneofDesc.Options().(*descriptorpb.FieldOptions), mir.E_EventType).(bool)
}

// Fields is a list of fields of a protobuf message.
type Fields []*Field

// FuncParamsPbTypes returns a list of field lowercase names followed by their pb types.
func (fs Fields) FuncParamsPbTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.FuncParamPbType() })
}

// FuncParamsMirTypes returns a list of field lowercase names followed by their mir types.
func (fs Fields) FuncParamsMirTypes() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return f.FuncParamMirType() })
}

// FuncParamsIDs returns a list of fields lowercase names as identifiers, without the types.
func (fs Fields) FuncParamsIDs() []jen.Code {
	return sliceutil.Transform(fs, func(i int, f *Field) jen.Code { return jen.Id(f.LowercaseName()) })
}
