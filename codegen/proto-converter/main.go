package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang/mock/mockgen/model"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use: "proto-converter",
		Short: "proto-converter – modifies the go files generated by protoc in order to " +
			"extend protobuf functionality for the use in Mir framework.",
		Args: cobra.ExactArgs(1),
		RunE: runConverter,
	}
)

func runConverter(cmd *cobra.Command, args []string) error {
	packageName := args[0]

	// Get the current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get the working directory: %v", err)
	}

	// Find the package with the .pb.go files.
	pkg, err := build.Import(packageName, wd, build.FindOnly)
	if err != nil {
		return fmt.Errorf("cannot find package: %v", packageName)
	}

	// Create a temporary folder for the generator program.
	genDir, err := ioutil.TempDir(pkg.Dir, "proto-converter-tmp-")
	if err != nil {
		return fmt.Errorf("error creating a temporary directory in %v: %v", pkg.Dir, err)
	}
	defer func() {
		err := os.RemoveAll(genDir)
		if err != nil {
			log.Printf("failed to remove the temporary directory %v: %v", genDir, err)
		}
	}()

	// Get the list of all structs (i.e., potential message types) in the package.
	structNames, err := getAllStructNamesInAPackage(packageName)
	if err != nil {
		return fmt.Errorf("failed to extract message types from package %v: %w", packageName, err)
	}

	// Create the generator program code.
	generatorProgram, err := writeGeneratorProgram(packageName, structNames)
	if err != nil {
		return err
	}

	runInDir(genDir)
	return nil
}

func writeGeneratorProgram(importPath string, structNames []string) ([]byte, error) {
	var program bytes.Buffer
	data := reflectData{
		ImportPath: importPath,
		Symbols:    symbols,
	}
	if err := generatorTemplate.Execute(&program, &data); err != nil {
		return nil, err
	}
	return program.Bytes(), nil
}

func runInDir(program []byte, dir string) (*model.Package, error) {
	// We use TempDir instead of TempFile so we can control the filename.
	tmpDir, err := ioutil.TempDir(dir, "gomock_reflect_")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("failed to remove temp directory: %s", err)
		}
	}()
	const progSource = "prog.go"
	var progBinary = "prog.bin"
	if runtime.GOOS == "windows" {
		// Windows won't execute a program unless it has a ".exe" suffix.
		progBinary += ".exe"
	}

	if err := ioutil.WriteFile(filepath.Join(tmpDir, progSource), program, 0600); err != nil {
		return nil, err
	}

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "build")
	if *buildFlags != "" {
		cmdArgs = append(cmdArgs, strings.Split(*buildFlags, " ")...)
	}
	cmdArgs = append(cmdArgs, "-o", progBinary, progSource)

	// Build the program.
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, buf)
	if err := cmd.Run(); err != nil {
		sErr := buf.String()
		if strings.Contains(sErr, `cannot find package "."`) &&
			strings.Contains(sErr, "github.com/golang/mock/mockgen/model") {
			fmt.Fprint(os.Stderr, "Please reference the steps in the README to fix this error:\n\thttps://github.com/golang/mock#reflect-vendoring-error.\n")
			return nil, err
		}
		return nil, err
	}

	return run(filepath.Join(tmpDir, progBinary))
}

// run the given program and parse the output as a model.Package.
func run(program string) (*model.Package, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}

	filename := f.Name()
	defer os.Remove(filename)
	if err := f.Close(); err != nil {
		return nil, err
	}

	// Run the program.
	cmd := exec.Command(program, "-output", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	f, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Process output.
	var pkg model.Package
	if err := gob.NewDecoder(f).Decode(&pkg); err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return &pkg, nil
}

func main() {
	err := cmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error executing the converter: %v\n", err)
		os.Exit(2)
	}
}

func getAllStructNamesInAPackage(sourcePkgName string) ([]string, error) {
	// Get the source dir of the package.
	pkg, err := importer.Default().Import(sourcePkgName)
	if err != nil {
		return nil, fmt.Errorf("could not obtain the source dir for package %v: %w", sourcePkgName, err)
	}
	sourceDir := pkg.Path()

	// Parse the source dir.
	fset := token.FileSet{}
	pkgs, err := parser.ParseDir(&fset, sourceDir, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse the package sources: %v", err)
	}

	// Get the ast of the target package.
	pkgAst, ok := pkgs[sourcePkgName]
	if !ok {
		return nil, fmt.Errorf("did not find package %v in %v", sourcePkgName, sourceDir)
	}

	var res []string
	for _, file := range pkgAst.Files {
		for _, decl := range file.Decls {
			decl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range decl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if typeSpec.Assign != token.NoPos {
					// Ignore type aliases.
					continue
				}

				res = append(res, typeSpec.Name.Name)
			}
		}
	}

	return res, nil
}

var generatorTemplate = template.Must(template.New("generator").Parse(`

`))
