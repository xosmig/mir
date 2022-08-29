package main

import (
	"flag"
	"fmt"
	"reflect"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/filecoin-project/mir/codegen/proto-converter/codegen/model"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/protogenutil"
	"github.com/filecoin-project/mir/pkg/util/reflectutil"
)

func generateReflectMethodsToListOneofOptions(plugin *protogen.Plugin, file *protogen.File) error {
	var g *protogen.GeneratedFile

	for _, msg := range file.Messages {
		if !model.ShouldGenerateMirType(msg.Desc) {
			// Skip structs that are not explicitly marked with Mir annotations.
			continue
		}

		for _, oneof := range msg.Oneofs {
			if g == nil {
				filename := fmt.Sprintf("%s.pb.mir.go", file.GeneratedFilenamePrefix)
				g = plugin.NewGeneratedFile(filename, file.GoImportPath)
				g.P("package ", file.GoPackageName)
				g.P()
			}

			interfaceName := g.QualifiedGoIdent(oneof.GoIdent)
			g.P("type ", interfaceName, " = ", "is", interfaceName)
			g.P()

			reflectType := g.QualifiedGoIdent(protogenutil.GoIdentByType(reflectutil.TypeOf[reflect.Type]()))
			g.P("func (*", msg.GoIdent, ") Reflect", oneof.GoName, "Options() []", reflectType, " {")
			g.P("\t", "return []", reflectType, " {")
			for _, field := range oneof.Fields {
				wrapperTypeName := g.QualifiedGoIdent(field.GoIdent)
				g.P("\t\t", "reflect.TypeOf((*", wrapperTypeName, ") (nil)),")
			}
			g.P("\t", "}")
			g.P("}")
			g.P()
		}
	}

	return nil
}

func main() {
	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		for _, f := range plugin.Files {
			err := generateReflectMethodsToListOneofOptions(plugin, f)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
