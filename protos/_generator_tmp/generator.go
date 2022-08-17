package main

import (
	"log"
	"os"

	"github.com/filecoin-project/mir/codegen/proto-converter/codegen"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb"
)

func main() {
	err := codegen.ProduceFile(
		"availabilitypb",
		"./pkg/pb/availabilitypb/structs/availabilitypb.mir.go",
		[]any{
			&availabilitypb.CertVerified{},
			&availabilitypb.RequestCertOrigin{},
		})

	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}
