package codegen

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/codegen/proto-converter/jenutil"
	"github.com/filecoin-project/mir/pkg/pb/mir"
)

func generateStruct(g *jen.File, structPointer any, structAst *ast.StructType) error {
	protoMsg, ok := structPointer.(protoreflect.ProtoMessage)
	if !ok {
		// This struct is not a ProtoMessage.
		return nil
	}

	msgDesc := protoMsg.ProtoReflect().Descriptor()
	msgOptions := msgDesc.Options().(*descriptorpb.MessageOptions)

	isMirEvent := proto.GetExtension(msgOptions, mir.E_Event).(bool)
	isMirStruct := proto.GetExtension(msgOptions, mir.E_Struct).(bool)
	if !isMirEvent && !isMirStruct {
		// This proto message is not a mir event.
		return nil
	}

	//for fieldID := 0; fieldID < msgDesc.Fields().Len(); fieldID++ {
	//	field := msgDesc.Fields().Get(fieldID)
	//	fmt.Println(field.Name())
	//	fmt.Println(field.FullName())
	//}

	structRefl := reflect.TypeOf(structPointer).Elem()

	var fields []jen.Code

	for i := 0; i < structRefl.NumField(); i++ {
		// Get go representation of the field.
		goField := structRefl.Field(i)
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
			return err
		}
		protoField := msgDesc.Fields().ByName(protoreflect.Name(protoName))

		// Get the desired type for the field.
		mirType, err := getMirType(goField, protoField)
		if err != nil {
			return err
		}

		fields = append(fields, jen.Id(goField.Name).Add(mirType))
	}

	g.Type().Id(structRefl.Name()).Struct(fields...)

	//convertedAst := &ast.StructType{
	//	Fields:
	//}
	//
	//structType := reflect.TypeOf(structPointer).Elem()
	//
	//
	//
	//
	//protoMessageInterface := reflect.TypeOf((*protoreflect.ProtoMessage)(nil)).Elem()
	//if !tp.Implements(protoMessageInterface) {
	//	// This struct is not a ProtoMessage
	//	return
	//}
	//
	//tp.

	return nil
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

func getMirType(goField reflect.StructField, protoField protoreflect.FieldDescriptor) (jen.Code, error) {
	protoFieldOptions := protoField.Options().(*descriptorpb.FieldOptions)
	mirType := proto.GetExtension(protoFieldOptions, mir.E_Type).(string)

	if mirType != "" {
		sepIdx := strings.LastIndex(mirType, ".")

		return jen.Qual(mirType[:sepIdx], mirType[sepIdx+1:]), nil
	}

	return jenutil.QualFromType(goField.Type), nil
}

func ProduceFile(packageName, outputPath string, structPointers []any) (err error) {
	g := jen.NewFile(packageName)

	for _, structPointer := range structPointers {
		err := generateStruct(g, structPointer, nil)
		if err != nil {
			return err
		}
	}

	// Open the output file.
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
		// Remove the output file in case of a failure to avoid causing compilation errors.
		if err != nil {
			_ = os.Remove(outputPath)
		}
	}()

	// Render the file.
	return g.Render(f)
}
