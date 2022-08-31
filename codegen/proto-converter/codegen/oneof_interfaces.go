package codegen

import (
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
func GenerateOneofInterfaces(
	inputDir, inputPackagePath string,
	msgs []*model.Message,
	parser *model.Parser,
) (err error) {

	jenFile := jen.NewFilePath(inputPackagePath)

	for _, msg := range msgs {
		fields, err := parser.ParseFields(msg)
		if err != nil {
			return err
		}

		for _, field := range fields {
			if oneof, ok := field.Type.(*model.Oneof); ok {
				err := generateOneofInterface(jenFile, oneof)
				if err != nil {
					return err
				}
			}
		}
	}

	return renderJenFile(jenFile, inputDir, "oneof_reflect.pb.mir.go", /*removeDirOnFail*/ false)
}
