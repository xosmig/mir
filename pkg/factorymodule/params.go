package factorymodule

import (
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/factorymodulepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

const (
	DefaultMsgBufSize = 1024 * 1024 // 1 MB
)

// ModuleGenerator is a function that generates a passive module from a module ID and the parameters.
// The parameters format is defined in .proto files (see file: protos/factorymodulepb/factorymodulepb.proto).
type ModuleGenerator func(id t.ModuleID, params *factorymodulepb.GeneratorParams) (modules.PassiveModule, error)

// ModuleParams contains the parameters for creating a factory module.
type ModuleParams struct {
	Generator  ModuleGenerator
	MsgBufSize int
}

// DefaultParams creates a valid ModuleParams with default values.
func DefaultParams(generator ModuleGenerator) ModuleParams {
	return ModuleParams{
		Generator:  generator,
		MsgBufSize: DefaultMsgBufSize,
	}
}
