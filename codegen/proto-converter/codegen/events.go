package codegen

import (
	"fmt"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/codegen/model"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/importerutil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

func EventsPackagePath(sourcePackagePath string) string {
	return sourcePackagePath + "/events"
}

func EventsPackageName(sourcePackagePath string) string {
	return sourcePackagePath[strings.LastIndex(sourcePackagePath, "/")+1:] + "events"
}

func EventsOutputDir(sourceDir string) string {
	return path.Join(sourceDir, "events")
}

func generateEventConstructorsRecursively(
	eventNode *model.EventNode,
	constructParent func(code jen.Code) *jen.Statement,
	eventRootType jen.Code,
	parentExtraFields model.Fields,
	jenFileBySourcePackagePath map[string]*jen.File,
	parser *model.Parser,
) error {

	fields, err := parser.ParseFields(eventNode.Message())
	if err != nil {
		return err
	}

	// Filter out the fields that are marked with option [(mir.omit_in_constructor) = true].
	fields = sliceutil.Filter(fields, func(_ int, field *model.Field) bool { return !field.OmitInConstructor() })

	// TODO: resolve potential name collisions with the parent fields.
	fieldsWithParent := append(parentExtraFields, fields...)

	// If this is an intermediate node in the hierarchy, recursively call the function for subtypes.
	if eventNode.IsEventClass() {
		for _, childNode := range eventNode.Children() {
			// constructThis is a function that takes the code to construct a child in the hierarchy
			// and constructs an event.
			constructThis := func(child jen.Code) *jen.Statement {
				return constructParent(
					eventNode.Message().NewMirType().ValuesFunc(func(group *jen.Group) {
						// Initialize other fields.
						for _, field := range fields {
							if field.Name != model.TypeOneofFieldName {
								group.Line().Id(field.Name).Op(":").Id(field.LowercaseName())
							}
						}

						// Initialize the Type field
						group.Line().Id(model.TypeOneofFieldName).Op(":").
							Add(childNode.OneofOption().ConstructMirWrapperType(child))

						// Put the closing bracket on a new line.
						group.Line()
					}),
				)
			}

			fieldsWithParentWithoutType := sliceutil.Filter(fieldsWithParent, func(i int, f *model.Field) bool {
				return f.Name != model.TypeOneofFieldName
			})

			err := generateEventConstructorsRecursively(
				/*eventNode*/ childNode,
				/*constructParent*/ constructThis,
				eventRootType,
				fieldsWithParentWithoutType,
				jenFileBySourcePackagePath,
				parser,
			)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// If this is an event (i.e., a leaf in the hierarchy), create the event constructor.

	// First, get a jen file to which the event constructor will be added.
	sourcePackage := eventNode.Message().PbPkgPath()
	jenFile, ok := jenFileBySourcePackagePath[sourcePackage]
	if !ok {
		jenFile = jen.NewFilePathName(EventsPackagePath(sourcePackage), EventsPackageName(sourcePackage))
		jenFileBySourcePackagePath[sourcePackage] = jenFile
	}

	// Generate the constructor.
	jenFile.Func().Id(eventNode.Message().Name()).Params(fieldsWithParent.FuncParamsMirTypes()...).Add(eventRootType).Block(
		jen.Return(constructParent(
			eventNode.Message().NewMirType().ValuesFunc(func(group *jen.Group) {
				for _, field := range fields {
					group.Line().Id(field.Name).Op(":").Id(field.LowercaseName())
				}
				group.Line()
			}),
		)),
	)

	return nil
}

// GenerateEventConstructors generates functions of form:
//
//     func [SrcPkg]events.[EventName]([EventParams]...) [RootEventType]
//
// TODO: add an example.
func GenerateEventConstructors(eventRoot *model.EventNode, parser *model.Parser) error {
	jenFileBySourcePackagePath := make(map[string]*jen.File)

	err := generateEventConstructorsRecursively(
		/*eventNode*/ eventRoot,
		/*constructParent*/ func(code jen.Code) *jen.Statement { return jen.Add(code) },
		/*eventRootType*/ eventRoot.Message().MirType(),
		/*parentExtraFields*/ nil,
		jenFileBySourcePackagePath,
		parser,
	)
	if err != nil {
		return fmt.Errorf("error generating event constructors: %w", err)
	}

	for sourcePackage, jenFile := range jenFileBySourcePackagePath {
		sourceDir, err := importerutil.GetSourceDirForPackage(sourcePackage)
		if err != nil {
			return fmt.Errorf("could not find the source directory for package %v: %w", sourcePackage, err)
		}

		outputDir := EventsOutputDir(sourceDir)
		err = renderJenFile(jenFile, outputDir, "events.mir.go", /*removeDirOnFail*/ true)
		if err != nil {
			return err
		}
	}

	return nil
}
