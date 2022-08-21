package codegen

import (
	"fmt"
	"os"
	"path"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/codegen/model"
)

func generateMirType(g *jen.File, msg *model.Message, oneofOptions []*model.OneofOption) error {
	if !msg.ShouldGenerateMirType() {
		// Ignore non-annotated messages.
		return nil
	}

	fields, err := msg.Fields(oneofOptions)
	if err != nil {
		return err
	}

	// Generate the struct.
	g.Type().Id(msg.Name()).StructFunc(func(group *jen.Group) {
		for _, field := range fields {
			group.Id(field.Name).Add(field.Type.MirType())
		}
	}).Line()

	// Generate the code for oneof fields.
	for _, field := range fields {
		if oneof, ok := field.Type.(*model.Oneof); ok {
			// Generate the interface.
			g.Type().Id(oneof.MirInterfaceName()).Interface(
				jen.Id(oneof.MirMethodName()).Params(),
				jen.Id("Pb").Params().Add(oneof.PbType()),
			).Line()

			// Generate the wrappers.
			for _, opt := range oneof.Options {
				g.Type().Id(opt.WrapperName).Struct(
					jen.Id(opt.Field.Name).Add(opt.Field.Type.MirType()),
				).Line()

				g.Func().Params(jen.Op("*").Id(opt.WrapperName)).Id(oneof.MirMethodName()).Params().Block().Line()

				g.Func().Params(jen.Id("w").Op("*").Id(opt.WrapperName)).Id("Pb").Params().Add(oneof.PbType()).Block(
					jen.Return(jen.Add(opt.NewPbWrapperType()).Values(
						jen.Id(opt.Field.Name).Op(":").Add(opt.Field.Type.ToPb(jen.Id("w").Dot(opt.Field.Name))),
					)),
				).Line()
			}
		}
	}

	// Generate New[Name] function.
	g.Func().Id("New" + msg.Name()).Params(fields.FuncParamsMirTypes()...).Add(msg.MirType()).Block(
		jen.Return().Add(msg.NewMirType()).ValuesFunc(func(group *jen.Group) {
			for _, field := range fields {
				group.Line().Id(field.Name).Op(":").Id(field.LowercaseName())
			}
			group.Line()
		}),
	).Line()

	// Generate Pb() method.
	g.Func().Params(jen.Id("m").Add(msg.MirType())).Id("Pb").Params().Add(msg.PbType()).Block(
		jen.Return().Add(msg.NewPbType()).ValuesFunc(func(group *jen.Group) {
			for _, field := range fields {
				group.Line().Id(field.Name).Op(":").Add(field.Type.ToPb(jen.Id("m").Dot(field.Name)))
			}
			group.Line()
		}),
	).Line()

	// Generate [Name]FromPb function.
	// NB: it would be nicer to generate .ToMir() methods for pb types, but this would cause a cyclic dependency.
	g.Func().Id(msg.Name() + "FromPb").Params(jen.Id("pb").Add(msg.PbType())).Add(msg.MirType()).Block(
		jen.Return().Add(msg.NewMirType()).ValuesFunc(func(group *jen.Group) {
			for _, field := range fields {
				group.Line().Id(field.Name).Op(":").Add(field.Type.ToMir(jen.Id("pb").Dot(field.Name)))
			}
			group.Line()
		}),
	).Line()

	return nil
}

func generateOneofOption(g *jen.File, opt *model.OneofOption) error {
	//g.Type().Id(opt.WrapperName).Struct(
	//	jen.Id(opt.FieldName).Add(opt.Type.MirType()),
	//).Line()
	//
	//g.Func().Params(jen.Op("*").Id(opt.WrapperName)).Id(opt.MethodName).Params().Block().Line()
	return nil
}

func GenerateMirTypes(inputDir, inputPackagePath string, msgs []*model.Message, oneofOptions []*model.OneofOption) (err error) {
	// Determine the output package and path.
	outputPackagePath := model.StructsPackagePath(inputPackagePath)
	outputDir := path.Join(inputDir, model.StructsPackageName)
	outputFile := path.Join(outputDir, model.StructsPackageName+".mir.go")

	g := jen.NewFilePath(outputPackagePath)

	// Generate Mir types for messages.
	for _, msg := range msgs {
		err := generateMirType(g, msg, oneofOptions)
		if err != nil {
			return err
		}
	}

	// Generate Mir types for oneof options.
	for _, opt := range oneofOptions {
		err := generateOneofOption(g, opt)
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
