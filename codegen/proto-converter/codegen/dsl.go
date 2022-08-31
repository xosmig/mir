package codegen

//import (
//	"fmt"
//	"strings"
//
//	"github.com/dave/jennifer/jen"
//
//	"github.com/filecoin-project/mir/codegen/proto-converter/codegen/model"
//	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
//	"github.com/filecoin-project/mir/pkg/dsl"
//	"github.com/filecoin-project/mir/pkg/pb/eventpb"
//	"github.com/filecoin-project/mir/pkg/util/reflectutil"
//	"github.com/filecoin-project/mir/pkg/util/sliceutil"
//)
//
////var eventRootType = jenutil.QualFromType(reflect.TypeOf(&eventpb.Event{}))
//
//func DslPackageName(pbPackagePath string) string {
//	return pbPackagePath[strings.LastIndex(pbPackagePath, "/")+1:] + "dsl"
//}
//
//func DslPackagePath(pbPackagePath string) string {
//	return fmt.Sprintf("%v/%v", pbPackagePath, DslPackageName(pbPackagePath))
//}
//
////func generateDslFunctionForEmittingEvent(g *jen.File, eventRoot *model.Message, oneofOptions []*model.OneofOption) error {
////
////	fields, err := eventRoot.Fields(oneofOptions)
////	if err != nil {
////		return err
////	}
////
////
////
////	//g.Func().Id(event.Name()).Params(fields.FuncParamsMirTypes()...).Add(eventRootType).Block(
////	//	jen.Return()
////	//	)
////
////	return nil
////}
//
//var (
//	dslPackagePath = reflectutil.TypeOf[dsl.Module]().PkgPath()
//	dslModule      = jen.Qual(dslPackagePath, "Module")
//	dslEmitEvent   = jen.Qual(dslPackagePath, "EmitEvent")
//	dslUponEvent   = jen.Qual(dslPackagePath, "UponEvent")
//	rootEventType  = reflectutil.TypeOf[*eventpb.Event]()
//)
//
//func generateDslFunctionsForEmittingEvent(
//	constructEvent func(code jen.Code) *jen.Statement,
//	parentFields model.Fields,
//	eventOpt *model.OneofOption,
//	fr *jenutil.FileRegistry,
//	parser *model.Parser,
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
//		for _, opt := range typeOneof.Options {
//			err := generateDslFunctionsForEmittingEvent(
//				/*emitParentEvent*/ emitThisEvent,
//				/*parentFields*/ fieldsWithParentWithoutType,
//				/*eventOpt*/ opt,
//				fr,
//				parser,
//			)
//
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//
//	// Generate a function for emitting the event.
//	outputPackage := DslPackagePath(msg.PbPkgPath())
//	g := fr.GetFile(outputPackage)
//
//	g.Func().Id(msg.Name()).Params(fieldsWithParent.FuncParamsMirTypes()...).Block(
//		dslEmitEvent.Params(constructEvent(eventOpt.NewMirWrapperType().Values(
//			jen.Line().Add(msg.Constructor()).Params(fields.FuncParamsIDs()...),
//			jen.Line(),
//		))),
//	).Line()
//
//	//g.Func().Id("Upon"+msg.Name()).Params(
//	//	jen.Id("m").Add(dslModule),
//	//	jen.Func().Id("handler").Params(fieldsWithParent.FuncParamsMirTypes()...).Id("error"),
//	//).Block(
//	//	uponParentEvent.Types(eventOpt.MirWrapperType()).Params()
//	//)
//
//	for _, field := range fields {
//		// Recursively call the generator on all subtypes.
//		if oneof, ok := field.Type.(*model.Oneof); ok && oneof.Name == "Type" {
//			fieldsWithParentWithoutType := sliceutil.Filter(fieldsWithParent, func(i int, f *model.Field) bool {
//				return f.Name != "Type"
//			})
//
//			emitThisEvent := jen.Qual(outputPackage, msg.Name())
//
//			for _, opt := range oneof.Options {
//				err := generateDslFunctionsForEmittingEvent(
//					/*emitParentEvent*/ emitThisEvent,
//					/*parentFields*/ fieldsWithParentWithoutType,
//					/*eventOpt*/ opt,
//					fr,
//					parser,
//				)
//
//				if err != nil {
//					return err
//				}
//			}
//		}
//	}
//
//	return nil
//}
//
//func generateDslFunctionsForEmittingEvents(
//	eventRoot *model.Message,
//	fr *jenutil.FileRegistry,
//	parser *model.Parser,
//) error {
//	fields, err := parser.ParseFields(eventRoot)
//	if err != nil {
//		return err
//	}
//
//	oneofFields := sliceutil.Filter(fields, func(i int, field *model.Field) bool {
//		_, ok := field.Type.(*model.Oneof)
//		return ok
//	})
//
//	if len(oneofFields) != 1 || oneofFields[0].Name != "Type" {
//		return fmt.Errorf("expected 1 oneof field named 'Type' in event root")
//	}
//
//	typeOneof := oneofFields[0].Type.(*model.Oneof)
//
//	for _, opt := range typeOneof.Options {
//		err := generateDslFunctionsForEmittingEvent(
//			/*emitParentEvent*/ emitThisEvent,
//			/*parentFields*/ fieldsWithParentWithoutType,
//			/*eventOpt*/ opt,
//			fr,
//			parser,
//		)
//
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func GenerateDslFunctions(inputDir, inputPackagePath string, msgs []*model.Message, parser *model.Parser) (err error) {
//	fr := jenutil.NewFileRegistry()
//
//	for _, msg := range msgs {
//		if msg.IsEventRoot() {
//
//		}
//	}
//
//	_ = outputFile
//	_ = g
//
//	return nil
//}
