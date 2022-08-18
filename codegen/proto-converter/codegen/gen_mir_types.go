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

func generateMirType(g *jen.File, pbGoStructPtrType reflect.Type) error {
	if _, ok := protoreflectutil.DescriptorForType(pbGoStructPtrType); !ok {
		// Ignore types that aren't protobuf messages.
		return nil
	}

	msg, err := MessageFromPbGoType(pbGoStructPtrType)
	if err != nil {
		return err
	}

	if !msg.ShouldGenerateMirType() {
		// Ignore non-annotated messages.
		return nil
	}

	fields, err := msg.Fields()
	if err != nil {
		return err
	}

	g.Type().Id(msg.Name).Struct(fields.StructParamsMirTypes()...)
	g.Line()

	// Generate FromPb function.
	structInit := jen.Dict{}
	for _, field := range fields {
		// Generates the code of form `Foo: MirType(protoMsg.Foo)`
		structInit[jen.Id(field.Name)] =
			jen.Add(field.FromPb).Params(jen.Id(msg.LowercaseName()).Dot(field.Name))
	}

	g.Func().Add(msg.FromPbFunc()).Params(msg.FuncParamPbType()).Add(msg.MirTypePtr()).
		Block(
			jen.Return().Op("&").Add(msg.MirType).Values(structInit),
		)
	g.Line()

	// Generate ToPb function.
	structInit = jen.Dict{}
	for _, field := range fields {
		// Generates the code of form `Foo: MirType(protoMsg.Foo)`
		structInit[jen.Id(field.Name)] =
			jen.Add(field.ToPb).Params(jen.Id(msg.LowercaseName()).Dot(field.Name))
	}

	g.Func().Add(msg.ToPbFunc()).Params(msg.FuncParamMirType()).Add(msg.PbTypePtr()).
		Block(
			jen.Return().Op("&").Add(msg.PbType).Values(structInit),
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

func GenerateMirTypes(inputDir string, pbGoStructPtrTypes []reflect.Type) (err error) {
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

	// Determine the output package and path.
	outputPackagePath := StructsPackagePath(inputPackagePath)
	outputDir := path.Join(inputDir, StructsPackageName)
	outputFile := path.Join(outputDir, StructsPackageName+".go")

	// Generate the code.
	g := jen.NewFilePath(outputPackagePath)
	for _, structPtrType := range pbGoStructPtrTypes {
		err := generateMirType(g, structPtrType)
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
