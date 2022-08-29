package codegen

import (
	"fmt"
	"os"
	"path"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/codegen/model"
)

func generateOneofInterface(g *jen.File, oneof *model.Oneof) error {
	// Generate [Msg]_[Oneof] interface.
	g.Type().Id(oneof.PbExportedInterfaceName()).Op("=").Id(oneof.PbNativeInterfaceName()).Line()

	// Generate [Msg]_[Oneof]Wrapper interface.
	g.Type().Id(oneof.PbExportedInterfaceName()+"Wrapper").Types(jen.Id("T").Any()).Interface(
		jen.Id(oneof.PbExportedInterfaceName()),
		jen.Id("Unwrap").Params().Op("*").Id("T"),
	).Line()

	// Generate Unwrap() method implementations
	for _, opt := range oneof.Options {
		g.Func().Params(jen.Id("w").Add(opt.PbWrapperType())).Id("Unwrap").Params().Add(opt.Field.Type.PbType()).Block(
			jen.Return(jen.Id("w").Dot(opt.Field.Name)),
		).Line()
	}

	return nil
}

// GenerateOneofInterfaces generates exported interfaces of the form "[Msg]_[Oneof]" "[Msg]_[Oneof]Wrapper" for all
// oneofs in the given messages, where [Msg] is the name of the message and [Oneof] is the name of the oneof.
func GenerateOneofInterfaces(inputDir, inputPackagePath string, msgs []*model.Message, oneofOptions []*model.OneofOption) (err error) {
	g := jen.NewFilePath(inputPackagePath)

	for _, msg := range msgs {
		fields, err := msg.Fields(oneofOptions)
		if err != nil {
			return err
		}

		for _, field := range fields {
			if oneof, ok := field.Type.(*model.Oneof); ok {
				err := generateOneofInterface(g, oneof)
				if err != nil {
					return err
				}
			}
		}
	}

	// Open the output file.
	outputFile := path.Join(inputDir, "pboneofs.mir.go")
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
