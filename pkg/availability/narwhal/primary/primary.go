package primary

import (
	"github.com/filecoin-project/mir/pkg/availability/narwhal/primary/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig = common.ModuleConfig

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:   "availability",
		Worker: "worker",
		Net:    "net",
		Crypto: "crypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
type ModuleParams = common.ModuleParams

func NewModule(mc *ModuleConfig, params *ModuleParams) modules.Module {
	m := dsl.NewModule(mc.Self)

}
