package codegen

import (
	"fmt"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/model/events"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/importerutil"
)

func DslPackagePath(sourcePackagePath string) string {
	return sourcePackagePath + "/dsl"
}

func DslPackageName(sourcePackagePath string) string {
	return sourcePackagePath[strings.LastIndex(sourcePackagePath, "/")+1:] + "dsl"
}

func DslOutputDir(sourceDir string) string {
	return path.Join(sourceDir, "dsl")
}

var (
	// Note: using reflection to determine this package path would cause a build dependency cycle.
	dslPackagePath = "github.com/filecoin-project/mir/pkg/dsl"
	dslModule      = jen.Qual(dslPackagePath, "Module")
	dslEmitEvent   = jen.Qual(dslPackagePath, "EmitEvent")
	dslUponEvent   = jen.Qual(dslPackagePath, "UponEvent")
)

func generateDslFunctionsForEmittingEventsRecursively(
	eventNode *events.EventNode,
	jenFileBySourcePackagePath map[string]*jen.File,
) {

	// If this is an internal node in the hierarchy, recursively call the function for subtypes.
	if eventNode.IsEventClass() {
		for _, child := range eventNode.Children() {
			generateDslFunctionsForEmittingEventsRecursively(child, jenFileBySourcePackagePath)
		}
		return
	}

	// Get a jen file to which the event constructor will be added.
	sourcePackage := eventNode.Message().PbPkgPath()
	jenFile, ok := jenFileBySourcePackagePath[sourcePackage]
	if !ok {
		jenFile = jen.NewFilePathName(DslPackagePath(sourcePackage), DslPackageName(sourcePackage))
		jenFileBySourcePackagePath[sourcePackage] = jenFile

		jenFile.Comment("Module-specific dsl functions for emitting events.")
		jenFile.Line()
	}

	// Generate the function for emitting the event
	funcParams := append(
		[]jen.Code{jen.Id("m").Add(dslModule)},
		eventNode.AllConstructorParameters().MirCode()...,
	)

	jenFile.Func().Id(eventNode.Name()).Params(funcParams...).Block(
		dslEmitEvent.Params(jen.Id("m"), eventNode.Constructor().Params(eventNode.AllConstructorParameters().IDs()...)),
	)
}

func generateDslFunctionsForHandlingEventsRecursively(
	eventNode *events.EventNode,
	uponParentEvent *jen.Statement,
	jenFileBySourcePackagePath map[string]*jen.File,
) {

	// Get a jen file to which the event constructor will be added.
	sourcePackage := eventNode.Message().PbPkgPath()
	jenFile, ok := jenFileBySourcePackagePath[sourcePackage]
	if !ok {
		jenFile = jen.NewFilePathName(DslPackagePath(sourcePackage), DslPackageName(sourcePackage))
		jenFileBySourcePackagePath[sourcePackage] = jenFile

		jenFile.Comment("Module-specific dsl functions for processing events.")
		jenFile.Line()
	}

	handlerParameters := eventNode.ThisNodeConstructorParameters()

	// Check if this is an internal node in the hierarchy.
	if eventNode.IsEventClass() {

		// Generate function for handling the event class.
		jenFile.Func().Id("Upon"+eventNode.Name()).Types(
			jen.Id("W").Id(eventNode.TypeOneof().MirWrapperInterfaceName()).Types(jen.Id("Ev")),
			jen.Id("Ev").Any(),
		).Params(
			jen.Id("m").Add(dslModule),
			jen.Id("handler").Func().Params(handlerParameters.MirCode()...).Id("error"),
		).Block(
			uponParentEvent.Types(eventNode.OneofOption().PbWrapperType()).Params(
				jen.Id("m"),
				jen.Func().Params(jen.Id("ev").Add(eventNode.Message().MirType())).Id("error").Block(
					jen.List(jen.Id("w"), jen.Id("ok")).Op(":=").
						Id("ev").Dot(eventNode.Parent().TypeOneof().Name).Op(".").Add(jen.Id("W")),
					jen.If(jen.Op("!").Id("ok")).Block(
						jen.Return(jen.Id("nil")),
					),
					jen.Return(jen.Id("handler").Params(jen.Id("w").Dot("Unwrap").Params())),
				),
			),
		)

		uponThisEvent := jen.Qual(DslPackagePath(eventNode.Message().PbPkgPath()), "Upon"+eventNode.Name())

		// Recursively invoke the function for the children in the hierarchy.
		for _, child := range eventNode.Children() {
			generateDslFunctionsForHandlingEventsRecursively(
				/*eventNode*/ child,
				/*uponParentEvent*/ uponThisEvent,
				jenFileBySourcePackagePath,
			)
		}
		return
	}

	// Generate the function for handling the event.
	jenFile.Func().Id("Upon"+eventNode.Name()).Params(
		jen.Id("m").Add(dslModule),
		jen.Id("handler").Func().Params(handlerParameters.MirCode()...).Id("error"),
	).Block(
		uponParentEvent.Types(eventNode.OneofOption().MirWrapperType()).Params(
			jen.Id("m"),
			jen.Func().Params(jen.Id("ev").Add(eventNode.Message().MirType()).Id("error")).Block(
				jen.Return(jen.Id("handler").Params(handlerParameters.IDs()...)),
			),
		),
	)
}

func GenerateDslFunctions(eventRoot *events.EventNode) error {
	jenFileBySourcePackagePath := make(map[string]*jen.File)

	generateDslFunctionsForEmittingEventsRecursively(
		/*eventNode*/ eventRoot,
		jenFileBySourcePackagePath,
	)

	for sourcePackage, jenFile := range jenFileBySourcePackagePath {
		sourceDir, err := importerutil.GetSourceDirForPackage(sourcePackage)
		if err != nil {
			return fmt.Errorf("could not find the source directory for package %v: %w", sourcePackage, err)
		}

		outputDir := DslOutputDir(sourceDir)
		err = renderJenFile(jenFile, outputDir, "dsl.mir.go" /*removeDirOnFail*/, true)
		if err != nil {
			return err
		}
	}

	return nil
}
