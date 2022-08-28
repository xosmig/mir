package codegen

//import (
//	"fmt"
//	"path"
//	"strings"
//
//	"github.com/dave/jennifer/jen"
//
//	"github.com/filecoin-project/mir/codegen/proto-converter/codegen/model"
//	"github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"
//	"github.com/filecoin-project/mir/pkg/dsl"
//	"github.com/filecoin-project/mir/pkg/util/reflectutil"
//	"github.com/filecoin-project/mir/pkg/util/sliceutil"
//)
//
////var eventRootType = jenutil.QualFromType(reflect.TypeOf(&eventpb.Event{}))
//
//func DslPackageName(pbPackagePath string) string {
//	return pbPackagePath[strings.LastIndex(pbPackagePath, "/")+1:] + "structs"
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
//	dslEmitEvent   = jen.Qual(dslPackagePath, "EmitEvent")
//	dslUponEvent   = jen.Qual(dslPackagePath, "UponEvent")
//)
//
//func generateDslFunctionsForEvent(
//	emitParentEvent *jen.Statement,
//	parentFields model.Fields,
//	eventOpt *model.OneofOption,
//	fr *jenutil.FileRegistry,
//	oneofOptions []*model.OneofOption,
//) error {
//	msg, ok := eventOpt.Field.Type.(*model.Message)
//	if !ok {
//		return fmt.Errorf("event %v is not a message", eventOpt.Field.Name)
//	}
//
//	fields, err := msg.Fields(oneofOptions)
//	if err != nil {
//		return err
//	}
//
//	fieldsWithParent := append(parentFields, fields...)
//
//	outputPackage := DslPackagePath(msg.PbPkgPath())
//	g := fr.GetFile(outputPackage)
//
//	g.Func().Id(msg.Name()).Params(fieldsWithParent.FuncParamsMirTypes()...).Block(
//		emitParentEvent.Params(
//			msg.Constructor().Params(fieldsWithParent.FuncParamsIDs()...),
//		),
//	).Line()
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
//				err := generateDslFunctionsForEvent(
//					/*emitParentEvent*/ emitThisEvent,
//					/*parentFields*/ fieldsWithParentWithoutType,
//					/*eventOpt*/ opt,
//					fr,
//					oneofOptions,
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
//func generateDslFunctions(eventRoot *model.Message, oneofOptions []*model.OneofOption) error {
//	fields, err := eventRoot.Fields(oneofOptions)
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
//		err := generateDslFunctionsForEvent(opt, oneofOptions)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func GenerateDslFunctions(inputDir, inputPackagePath string, msgs []*model.Message, oneofOptions []*model.OneofOption) (err error) {
//	// Determine the output package and path.
//	outputPackagePath := DslPackagePath(inputPackagePath)
//	outputDir := path.Join(inputDir, DslPackageName(inputPackagePath))
//	outputFile := path.Join(outputDir, DslPackageName(inputPackagePath)+".mir.go")
//
//	g := jen.NewFilePath(outputPackagePath)
//
//	events := sliceutil.Filter(msgs, func(i int, msg *model.Message) bool { return msg.IsMirEvent() })
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
