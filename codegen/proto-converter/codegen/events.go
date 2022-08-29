package codegen

//import (
//	"fmt"
//
//	"github.com/dave/jennifer/jen"
//
//	"github.com/filecoin-project/mir/codegen/proto-converter/codegen/model"
//	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
//	"github.com/filecoin-project/mir/pkg/util/sliceutil"
//)
//
//func generateEventConstructorsRec(
//	constructParent func(code jen.Code) *jen.Statement,
//	parentFields model.Fields,
//	eventOpt *model.OneofOption,
//	fr *jenutil.FileRegistry,
//	parser *model.Parser,
//	eventRootType jen.Code,
//) error {
//	msg, ok := eventOpt.Field.Type.(*model.Message)
//	if !ok {
//		return fmt.Errorf("event %v is not a message", eventOpt.Field.Name)
//	}
//
//	fields, err := parser.ParseFields(msg)
//	if err != nil {
//		return err
//	}
//
//	// TODO: resolve potential name collisions with the parent fields.
//	fieldsWithParent := append(parentFields, fields...)
//
//	// If this is an intermediate node in the hierarchy, recursively call the function for subtypes.
//	if typeOneof, ok := getTypeOneof(fields); ok {
//		constructThis := func(child jen.Code) *jen.Statement {
//			return constructParent(
//				eventOpt.ConstructMirWrapperType(
//					msg.NewMirType().ValuesFunc(func(group *jen.Group) {
//						for _, field := range fields {
//							group.Line().Id(field.Name).Op(":").Id(field.LowercaseName())
//						}
//						group.Line().Id(typeOneof.Name).Op(":").Add(child)
//						group.Line()
//					}),
//				),
//			)
//		}
//
//		fieldsWithParentWithoutType := sliceutil.Filter(fieldsWithParent, func(i int, f *model.Field) bool {
//			return f.Name != typeOneof.Name
//		})
//
//		for _, opt := range typeOneof.Options {
//			err := generateEventConstructorsRec(
//				/*constructParent*/ constructThis,
//				/*parentFields*/ fieldsWithParentWithoutType,
//				/*eventOpt*/ opt,
//				fr,
//				parser,
//				eventRootType,
//			)
//
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//
//	// If this is a leaf node in the hierarchy, create the event constructor.
//	outputPackage := DslPackagePath(msg.PbPkgPath())
//	g := fr.GetFile(outputPackage)
//
//	g.Func().Id(msg.Name()).Params(fieldsWithParent.FuncParamsMirTypes()...).Add(eventRootType).Block(
//		jen.Return(constructParent(
//			eventOpt.ConstructMirWrapperType(
//				msg.NewMirType().ValuesFunc(func(group *jen.Group) {
//					for _, field := range fields {
//						group.Line().Id(field.Name).Op(":").Id(field.LowercaseName())
//					}
//					group.Line()
//				}),
//			),
//		)),
//	)
//
//	return nil
//}
//
//// GenerateEventConstructors generates functions of form:
////
////     func [SrcPkg]events.[EventName]([EventParams]...) [RootEventType]
////
//// TODO: add an example.
////
//// TODO: generalize event hierarchy traversal.
//func GenerateEventConstructors(eventRoot *model.Message, parser *model.Parser) error {
//	fields, err := parser.ParseFields(eventRoot)
//	if err != nil {
//		return err
//	}
//
//	typeOneof, ok := getTypeOneof(fields)
//	if !ok {
//		return fmt.Errorf("event root message should have a oneof named 'type' or 'Type'")
//	}
//
//	for _, opt := range typeOneof.Options {
//		err := generateEventConstructorsRec(
//			/*constructParent*/ emitThisEvent,
//			/*parentFields*/ fieldsWithParentWithoutType,
//			/*eventOpt*/ opt,
//			fr,
//			parser,
//			/*eventRootType*/ eventRoot.MirType(),
//		)
//
//		if err != nil {
//			return err
//		}
//	}
//}
//
//func getTypeOneof(fields model.Fields) (*model.Oneof, bool) {
//	for _, field := range fields {
//		// Recursively call the generator on all subtypes.
//		if oneof, ok := field.Type.(*model.Oneof); ok && oneof.Name == "Type" {
//			return oneof, true
//		}
//	}
//	return nil, false
//}
