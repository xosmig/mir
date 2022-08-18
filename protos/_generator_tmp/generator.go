package main

import (
	"log"
	"os"
	"reflect"

	"github.com/filecoin-project/mir/codegen/proto-converter/codegen"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb"
)

func main() {
	err := codegen.GenerateMirTypes(
		"pkg/pb/availabilitypb",
		[]reflect.Type{
			reflect.TypeOf(&availabilitypb.CertVerified{}),
			reflect.TypeOf(&availabilitypb.RequestCertOrigin{}),
		})

	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}
