package types

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/protoreflectutil"
	"github.com/filecoin-project/mir/pkg/pb/mir"
)

type Parser struct {
	msgCache    map[reflect.Type]*Message
	fieldsCache map[reflect.Type]Fields
}

func NewParser() *Parser {
	return &Parser{
		msgCache:    make(map[reflect.Type]*Message),
		fieldsCache: make(map[reflect.Type]Fields),
	}
}

func (p *Parser) ParseMessages(pbGoStructPtrTypes []reflect.Type) ([]*Message, error) {
	var msgs []*Message

	for _, ptrType := range pbGoStructPtrTypes {
		if ptrType.Kind() != reflect.Pointer || ptrType.Elem().Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected a pointer to a struct, got %v", ptrType)
		}

		// Parse messages.
		if protoreflectutil.IsProtoMessage(ptrType) {
			msg, err := p.ParseMessage(ptrType)
			if err != nil {
				return nil, err
			}

			msgs = append(msgs, msg)
			continue
		}
	}

	return msgs, nil
}

// ParseMessage returns the message corresponding to the given protobuf-generated struct type.
func (p *Parser) ParseMessage(pbGoStructPtr reflect.Type) (*Message, error) {
	if tp, ok := p.msgCache[pbGoStructPtr]; ok {
		return tp, nil
	}

	protoDesc, ok := protoreflectutil.DescriptorForType(pbGoStructPtr)
	if !ok {
		return nil, fmt.Errorf("%T is not a protobuf message", pbGoStructPtr)
	}

	pbStructType := jenutil.QualFromType(pbGoStructPtr.Elem())

	shouldGenerateMirType := ShouldGenerateMirType(protoDesc)

	var pkgPath string
	var mirStructType jen.Code

	if shouldGenerateMirType {
		// The type of the struct that will be generated.
		pkgPath = StructsPackagePath(pbGoStructPtr.Elem().PkgPath())
		mirStructType = jen.Qual(pkgPath, pbGoStructPtr.Elem().Name())
	} else {
		// The original type generated by protoc.
		pkgPath = pbGoStructPtr.Elem().PkgPath()
		mirStructType = pbStructType
	}

	msg := &Message{
		shouldGenerateMirType: shouldGenerateMirType,
		pbStructType:          pbStructType,
		mirStructType:         mirStructType,
		protoDesc:             protoDesc,
		pbGoStructPtrReflect:  pbGoStructPtr,
	}

	p.msgCache[pbGoStructPtr] = msg
	return msg, nil
}

func (p *Parser) ParseOneofOption(message *Message, ptrType reflect.Type) (*OneofOption, error) {
	if !protoreflectutil.IsOneofOption(ptrType) {
		return nil, fmt.Errorf("%v is not a oneof option", ptrType)
	}

	// Get the go representation of the field.
	if ptrType.Elem().NumField() != 1 {
		return nil, fmt.Errorf("protoc-generated oneof wrapper must have exactly 1 exported field")
	}
	goField := ptrType.Elem().Field(0)

	// Get the protobuf representation of the field
	protoName, err := getProtoNameOfField(goField)
	if err != nil {
		return nil, fmt.Errorf("error parsing the name of proto field in oneof wrapper %v: %w", ptrType.Elem(), err)
	}
	protoField := message.protoDesc.Fields().ByName(protoName)

	// Parse the field information.
	field, err := p.parseField(goField, protoField)
	if err != nil {
		return nil, fmt.Errorf("error parsing oneof option %v: %w", ptrType.Name(), err)
	}

	// Return the resulting oneof option.
	return &OneofOption{
		PbWrapperReflect: ptrType,
		WrapperName:      ptrType.Elem().Name(),
		Field:            field,
	}, nil
}

// ParseFields parses the fields of a message for which a Mir type is being generated.
func (p *Parser) ParseFields(m *Message) (Fields, error) {
	// Return the cached value if present.
	if cachedFields, ok := p.fieldsCache[m.pbGoStructPtrReflect]; ok {
		return cachedFields, nil
	}

	if !m.ShouldGenerateMirType() {
		return nil, fmt.Errorf("fields can only be parsed for messages with Mir-generated types")
	}

	var fields Fields

	for i := 0; i < m.pbGoStructPtrReflect.Elem().NumField(); i++ {
		// Get go representation of the field.
		goField := m.pbGoStructPtrReflect.Elem().Field(i)
		if !ast.IsExported(goField.Name) {
			// Skip unexported fields.
			continue
		}

		// Process oneof fields.
		if oneofProtoName, ok := goField.Tag.Lookup("protobuf_oneof"); ok {
			oneofOptionsGoTypes := reflect.Zero(m.pbGoStructPtrReflect).
				MethodByName("Reflect" + goField.Name + "Options").Call([]reflect.Value{})[0].
				Interface().([]reflect.Type)

			var options []*OneofOption
			for _, optionGoType := range oneofOptionsGoTypes {
				opt, err := p.ParseOneofOption(m, optionGoType)
				if err != nil {
					return nil, err
				}

				options = append(options, opt)
			}

			fields = append(fields, &Field{
				Name: goField.Name,
				Type: &Oneof{
					Name:    goField.Name,
					Parent:  m,
					Options: options,
				},
				ProtoDesc: m.protoDesc.Oneofs().ByName(protoreflect.Name(oneofProtoName)),
			})
			continue
		}

		// Get protobuf representation of the field.
		protoName, err := getProtoNameOfField(goField)
		if err != nil {
			return nil, err
		}
		protoField := m.protoDesc.Fields().ByName(protoName)

		// Create the Field struct.
		field, err := p.parseField(goField, protoField)
		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	p.fieldsCache[m.pbGoStructPtrReflect] = fields
	return fields, nil
}

// parseField extracts the information about the field necessary for code generation.
func (p *Parser) parseField(goField reflect.StructField, protoField protoreflect.FieldDescriptor) (*Field, error) {
	tp, err := p.getFieldType(goField.Type, protoField)
	if err != nil {
		return nil, err
	}

	return &Field{
		Name:      goField.Name,
		Type:      tp,
		ProtoDesc: protoField,
	}, nil
}

func (p *Parser) getFieldType(goType reflect.Type, protoField protoreflect.FieldDescriptor) (Type, error) {
	// TODO: Since maps are not currently used, I didn't bother supporting them yet.
	if goType.Kind() == reflect.Map {
		return nil, fmt.Errorf("map fields are not supported yet")
	}

	// Check if the field is repeated.
	if goType.Kind() == reflect.Slice {
		underlying, err := p.getFieldType(goType.Elem(), protoField)
		if err != nil {
			return nil, err
		}
		return Slice{underlying}, nil
	}

	// Check if the field has (mir.type) option specified.
	protoFieldOptions := protoField.Options().(*descriptorpb.FieldOptions)
	mirTypeOption := proto.GetExtension(protoFieldOptions, mir.E_Type).(string)
	if mirTypeOption != "" {
		sepIdx := strings.LastIndex(mirTypeOption, ".")
		return Castable{
			pbType:  jenutil.QualFromType(goType),
			mirType: jen.Qual(mirTypeOption[:sepIdx], mirTypeOption[sepIdx+1:]),
		}, nil
	}

	// Check if the field is a message.
	if protoreflectutil.IsProtoMessage(goType) {
		msg, err := p.ParseMessage(goType)
		if err != nil {
			return nil, err
		}

		return msg, nil
	}

	return Same{jenutil.QualFromType(goType)}, nil
}