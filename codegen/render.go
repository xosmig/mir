package codegen

import (
	"fmt"
	"os"
	"path"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/util/buildutil"
)

func RenderJenFile(jenFile *jen.File, outputDir, outputFileName string) (err error) {
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

func RenderJenFiles(
	jenFileBySourcePackagePath map[string]*jen.File,
	outputDirBySourceDir func(string) string,
	outputFileName string,
) error {
	for sourcePackage, jenFile := range jenFileBySourcePackagePath {
		sourceDir, err := buildutil.GetSourceDirForPackage(sourcePackage)
		if err != nil {
			return fmt.Errorf("could not find the source directory for package %v: %w", sourcePackage, err)
		}

		outputDir := outputDirBySourceDir(sourceDir)
		err = RenderJenFile(jenFile, outputDir, outputFileName)
		if err != nil {
			return err
		}
	}

	return nil
}
