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

	// Check if this is an internal node in the hierarchy.
	if eventNode.IsEventClass() {
		jenFile.Func().Id("Upon"+eventNode.Name()).Types(
			jen.Id("ChildWrapper").Id(eventNode.TypeOneof().MirWrapperInterfaceName()).Params(jen.Id("Child")),
			jen.Id("Child").Any(),
		).Params(jen.Id("m").Add(dslModule), jen.Id("handler").Func().Params())

		for _, child := range eventNode.Children() {
			generateDslFunctionsForEmittingEventsRecursively(child, jenFileBySourcePackagePath)
		}

		// // UponEvent registers a handler for the given availability layer event type.
		//func UponEvent[EvWrapper apb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
		//	dsl.UponEvent[*eventpb.Event_Availability](m, func(ev *apb.Event) error {
		//		evWrapper, ok := ev.Type.(EvWrapper)
		//		if !ok {
		//			return nil
		//		}
		//		return handler(evWrapper.Unwrap())
		//	})
		//}
		return
	}

	// Generate the function for emitting the event
	funcParams := append(
		[]jen.Code{jen.Id("m").Add(dslModule)},
		eventNode.AllConstructorParameters().MirCode()...,
	)

	// TODO
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
