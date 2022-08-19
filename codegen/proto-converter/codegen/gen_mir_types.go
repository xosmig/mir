package codegen

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/util/protoreflectutil"
)

const StructsPackageName = "structs"

func StructsPackagePath(pbPackagePath string) string {
	return pbPackagePath + "/" + StructsPackageName
}

func generateMirType(g *jen.File, msg *Message) error {
	if !msg.ShouldGenerateMirType() {
		// Ignore non-annotated messages.
		return nil
	}

	fields, err := msg.Fields()
	if err != nil {
		return err
	}

	g.Type().Id(msg.Name()).Struct(fields.StructParamsMirTypes()...).Line()

	// Generate ToPb method.
	g.Func().Params(jen.Id("m").Add(msg.MirType())).Id("ToPb").Params().Add(msg.PbType()).Block(
		jen.Return().Add(msg.NewPbType()).ValuesFunc(func(group *jen.Group) {
			for _, field := range fields {
				group.Id(field.Name).Op(":").Add(field.Type.ToPb(jen.Id("m").Dot(field.Name)))
			}
		}),
	)

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

func GenerateMirTypes(inputDir, inputPackagePath string, msgs []*Message) (err error) {
	// Determine the output package and path.
	outputPackagePath := StructsPackagePath(inputPackagePath)
	outputDir := path.Join(inputDir, StructsPackageName)
	outputFile := path.Join(outputDir, StructsPackageName+".mir.go")

	// Generate the code.
	g := jen.NewFilePath(outputPackagePath)
	for _, msg := range msgs {
		err := generateMirType(g, msg)
		if err != nil {
			return err
		}
	}

	// Create the directory if needed.
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Open the output file.
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}

	defer func() {
		_ = f.Close()
		// Remove the output directory in case of a failure to avoid causing compilation errors.
		if err != nil {
			_ = os.RemoveAll(outputDir)
		}
	}()

	// Render the file.
	return g.Render(f)
}

func generateToMirMethod(g *jen.File, msg *Message) error {
	if !msg.ShouldGenerateMirType() {
		// Ignore non-annotated messages.
		return nil
	}

	fields, err := msg.Fields()
	if err != nil {
		return err
	}

	g.Func().Params(jen.Id("pb").Add(msg.PbType())).Id("ToMir").Params().Add(msg.MirType()).Block(
		jen.Return().Add(msg.NewMirType()).ValuesFunc(func(group *jen.Group) {
			for _, field := range fields {
				group.Id(field.Name).Op(":").Add(field.Type.ToMir(jen.Id("pb").Dot(field.Name)))
			}
		}),
	)

	return nil
}

func GenerateToMirMethods(inputDir string, inputPackagePath string, msgs []*Message) error {
	// Generate the code.
	g := jen.NewFilePath(inputPackagePath)
	for _, msg := range msgs {
		err := generateToMirMethod(g, msg)
		if err != nil {
			return err
		}
	}

	// Open the output file.
	outputFile := path.Join(inputDir, "to_mir.mir.go")
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}

	defer func() {
		_ = f.Close()
		// Remove the output file in case of a failure to avoid causing compilation errors.
		if err != nil {
			_ = os.Remove(outputFile)
		}
	}()

	// Render the file.
	return g.Render(f)
}

func GetMessagesFromGoTypes(pbGoStructPtrTypes []reflect.Type) ([]*Message, error) {
	var msgs []*Message
	for _, ptrType := range pbGoStructPtrTypes {
		if _, ok := protoreflectutil.DescriptorForType(ptrType); !ok {
			// Ignore types that aren't protobuf messages.
			continue
		}

		msg, err := MessageFromPbGoType(ptrType)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func GenerateAll(inputDir string, pbGoStructPtrTypes []reflect.Type) error {
	if len(pbGoStructPtrTypes) == 0 {
		// Nothing to do.
		return nil
	}

	// Determine the input package.
	inputPackagePath := pbGoStructPtrTypes[0].Elem().PkgPath()

	// Check that all structs are in the same package.
	for _, ptrType := range pbGoStructPtrTypes {
		if ptrType.Elem().PkgPath() != inputPackagePath {
			return fmt.Errorf("passed structs are in different packages: %v and %v",
				inputPackagePath, ptrType.Elem().PkgPath())
		}
	}

	msgs, err := GetMessagesFromGoTypes(pbGoStructPtrTypes)
	if err != nil {
		return err
	}

	err = GenerateMirTypes(inputDir, inputPackagePath, msgs)
	if err != nil {
		return err
	}

	err = GenerateToMirMethods(inputDir, inputPackagePath, msgs)
	if err != nil {
		return err
	}

	return nil
}
