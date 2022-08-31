package codegen

import (
	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/model/types"
)

func generateMirType(g *jen.File, msg *types.Message, parser *types.Parser) error {
	if !msg.ShouldGenerateMirType() {
		// Ignore non-annotated messages.
		return nil
	}

	fields, err := parser.ParseFields(msg)
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
		if oneof, ok := field.Type.(*types.Oneof); ok {
			// Generate the interface.
			g.Type().Id(oneof.MirInterfaceName()).Interface(
				jen.Id(oneof.MirMethodName()).Params(),
				jen.Id("Pb").Params().Add(oneof.PbType()),
			).Line()

			g.Func().Id(oneof.MirInterfaceName()+"FromPb").Params(jen.Id("pb").Add(oneof.PbType())).Add(oneof.MirType()).Block(
				jen.Switch(jen.Id("pb").Op(":=").Id("pb").Dot("(type)")).BlockFunc(func(group *jen.Group) {
					for _, opt := range oneof.Options {
						group.Case(opt.PbWrapperType()).Block(
							jen.Return(jen.Add(opt.NewMirWrapperType()).Values(
								jen.Id(opt.Field.Name).Op(":").Add(opt.Field.Type.ToMir(jen.Id("pb").Dot(opt.Field.Name))),
							)),
						)
					}
				}),
				jen.Return(jen.Nil()),
			).Line()

			// Generate the wrappers.
			for _, opt := range oneof.Options {
				g.Type().Id(opt.WrapperName).Struct(
					jen.Id(opt.Field.Name).Add(opt.Field.Type.MirType()),
				).Line()

				g.Func().Params(opt.MirWrapperType()).Id(oneof.MirMethodName()).Params().Block().Line()

				g.Func().Params(jen.Id("w").Add(opt.MirWrapperType())).Id("Pb").Params().Add(oneof.PbType()).Block(
					jen.Return(jen.Add(opt.NewPbWrapperType()).Values(
						jen.Id(opt.Field.Name).Op(":").Add(opt.Field.Type.ToPb(jen.Id("w").Dot(opt.Field.Name))),
					)),
				).Line()
			}
		}
	}

	//// Generate New[Name] function.
	//g.Func().Id(msg.ConstructorName()).Params(fields.FuncParamsMirTypes()...).Add(msg.MirType()).Block(
	//	jen.Return().Add(msg.NewMirType()).ValuesFunc(func(group *jen.Group) {
	//		for _, field := range fields {
	//			group.Line().Id(field.Name).Op(":").Id(field.LowercaseName())
	//		}
	//		group.Line()
	//	}),
	//).Line()

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

func GenerateMirTypes(inputDir, sourcePackagePath string, msgs []*types.Message, parser *types.Parser) (err error) {
	jenFile := jen.NewFilePathName(
		types.StructsPackagePath(sourcePackagePath),
		types.StructsPackageName(sourcePackagePath),
	)

	// Generate Mir types for messages.
	for _, msg := range msgs {
		err := generateMirType(jenFile, msg, parser)
		if err != nil {
			return err
		}
	}

	return renderJenFile(jenFile, types.StructsOutputDir(inputDir), "structs.mir.go" /*removeDirOnFail*/, true)
}
