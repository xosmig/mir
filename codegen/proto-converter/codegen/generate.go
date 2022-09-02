package codegen

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/model/events"
	"github.com/filecoin-project/mir/codegen/proto-converter/model/types"
	"github.com/filecoin-project/mir/codegen/proto-converter/util/importerutil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// GenerateAll is used by the generator program (which, in turn, is generated by proto-converter).
func GenerateAll(pbGoStructPtrTypes []reflect.Type) error {
	if len(pbGoStructPtrTypes) == 0 {
		// Nothing to do.
		return nil
	}

	// Determine the input package.
	inputPackagePath := pbGoStructPtrTypes[0].Elem().PkgPath()

	// Get the directory with input sources.
	inputDir, err := importerutil.GetSourceDirForPackage(inputPackagePath)
	if err != nil {
		return err
	}

	// Check that all structs are in the same package.
	for _, ptrType := range pbGoStructPtrTypes {
		if ptrType.Elem().PkgPath() != inputPackagePath {
			return fmt.Errorf("passed structs are in different packages: %v and %v",
				inputPackagePath, ptrType.Elem().PkgPath())
		}
	}

	parser := types.NewParser()
	msgs, err := parser.ParseMessages(pbGoStructPtrTypes)
	if err != nil {
		return err
	}

	err = GenerateOneofInterfaces(inputDir, inputPackagePath, msgs, parser)
	if err != nil {
		return err
	}

	err = GenerateMirTypes(inputDir, inputPackagePath, msgs, parser)
	if err != nil {
		return err
	}

	// Look for the root of the event hierarchy.
	eventRootMessages := sliceutil.Filter(msgs, func(_ int, msg *types.Message) bool { return msg.IsEventRoot() })

	for _, eventRootMessage := range eventRootMessages {
		eventParser := events.NewParser(parser)

		eventRoot, err := eventParser.ParseEventHierarchy(eventRootMessage)
		if err != nil {
			return err
		}

		err = GenerateEventConstructors(eventRoot)
		if err != nil {
			return fmt.Errorf("error generating event constructors: %w", err)
		}

		err = GenerateDslFunctions(eventRoot)
		if err != nil {
			return fmt.Errorf("error generating dsl functions: %w", err)
		}
	}

	return nil
}

func renderJenFile(jenFile *jen.File, outputDir, outputFileName string) (err error) {
	// Create the directory if needed.
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Open the output file.
	filePath := path.Join(outputDir, outputFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}

	defer func() {
		_ = file.Close()
		// Remove the output file in case of a failure to avoid causing compilation errors.
		if err != nil {
			_ = os.Remove(filePath)
		}
	}()

	// Render the file.
	return jenFile.Render(file)
}

func renderJenFiles(
	jenFileBySourcePackagePath map[string]*jen.File,
	outputDirBySourceDir func(string) string,
	outputFileName string,
) error {
	for sourcePackage, jenFile := range jenFileBySourcePackagePath {
		sourceDir, err := importerutil.GetSourceDirForPackage(sourcePackage)
		if err != nil {
			return fmt.Errorf("could not find the source directory for package %v: %w", sourcePackage, err)
		}

		outputDir := outputDirBySourceDir(sourceDir)
		err = renderJenFile(jenFile, outputDir, outputFileName)
		if err != nil {
			return err
		}
	}

	return nil
}
