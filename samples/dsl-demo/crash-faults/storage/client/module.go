package client

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids.
type ModuleConfig struct {
	Self t.ModuleID // id of this module
	Net  t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self: "availability",
		Net:  "net",
	}
}

func NewModule(mc *ModuleConfig) {
	m := dsl.NewModule(mc.Self)

	dsl.UponNewRequests(m, func(requests []*requestpb.Request) error {
		for _, req := range requests {
			dsl.SendMessage(mc.Net)
		}
	})
}
