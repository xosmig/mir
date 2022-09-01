package codegen

import (
	"fmt"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/model/events"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/importerutil"
)

func generateEventConstructorsRecursively(
	eventNode *events.EventNode,
	constructParent func(code jen.Code) *jen.Statement,
	eventRootType jen.Code,
	jenFileBySourcePackagePath map[string]*jen.File,
) {

	// If this is an intermediate node in the hierarchy, recursively call the function for subtypes.
	if eventNode.IsEventClass() {
		for _, childNode := range eventNode.Children() {
			// constructThis is a function that takes the code to construct a child in the hierarchy
			// and constructs an event.
			constructThis := func(child jen.Code) *jen.Statement {
				return constructParent(
					eventNode.Message().NewMirType().ValuesFunc(func(group *jen.Group) {
						// Initialize fields other than the event type.
						for _, param := range eventNode.ThisNodeConstructorParameters().Slice {
							group.Line().Id(param.Field.Name).Op(":").Id(param.ParamName)
						}

						// Initialize the Type field
						group.Line().Id(eventNode.TypeOneof().Name).Op(":").
							Add(childNode.OneofOption().ConstructMirWrapperType(child))

						// Put the closing bracket on a new line.
						group.Line()
					}),
				)
			}

			generateEventConstructorsRecursively(
				/*eventNode*/ childNode,
				/*constructParent*/ constructThis,
				eventRootType,
				jenFileBySourcePackagePath,
			)
		}

		return
	}

	// If this is an event (i.e., a leaf in the hierarchy), create the event constructor.

	// Get a jen file to which the event constructor will be added.
	sourcePackage := eventNode.Message().PbPkgPath()
	jenFile, ok := jenFileBySourcePackagePath[sourcePackage]
	if !ok {
		jenFile = jen.NewFilePathName(events.PackagePath(sourcePackage), events.PackageName(sourcePackage))
		jenFileBySourcePackagePath[sourcePackage] = jenFile
	}

	// Generate the constructor.
	jenFile.Func().Id(eventNode.Name()).Params(
		eventNode.AllConstructorParameters().MirCode()...,
	).Add(eventRootType).Block(
		jen.Return(constructParent(
			eventNode.Message().NewMirType().ValuesFunc(func(group *jen.Group) {
				for _, param := range eventNode.ThisNodeConstructorParameters().Slice {
					group.Line().Id(param.Field.Name).Op(":").Id(param.ParamName)
				}
				group.Line()
			}),
		)),
	)
}

// GenerateEventConstructors generates functions of form:
//
//     func [SrcPkg]events.[EventName]([EventParams]...) [RootEventType]
//
// TODO: add an example.
func GenerateEventConstructors(eventRoot *events.EventNode) error {
	jenFileBySourcePackagePath := make(map[string]*jen.File)

	constructPbRootEventFromMirRootEvent := func(mirRootEvent jen.Code) *jen.Statement {
		return eventRoot.Message().ToPb(mirRootEvent)
	}

	generateEventConstructorsRecursively(
		/*eventNode*/ eventRoot,
		/*constructParent*/ constructPbRootEventFromMirRootEvent,
		/*eventRootType*/ eventRoot.Message().PbType(),
		jenFileBySourcePackagePath,
	)

	for sourcePackage, jenFile := range jenFileBySourcePackagePath {
		sourceDir, err := importerutil.GetSourceDirForPackage(sourcePackage)
		if err != nil {
			return fmt.Errorf("could not find the source directory for package %v: %w", sourcePackage, err)
		}

		outputDir := events.OutputDir(sourceDir)
		err = renderJenFile(jenFile, outputDir, "events.mir.go" /*removeDirOnFail*/, true)
		if err != nil {
			return err
		}
	}

	return nil
}