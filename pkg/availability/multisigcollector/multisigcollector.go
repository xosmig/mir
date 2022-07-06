package multisigcollector

import (
	adsl "github.com/filecoin-project/mir/pkg/availability/availabilitydsl"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/mempool/mempooldsl"
	"github.com/filecoin-project/mir/pkg/modules"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig struct {
	Self    t.ModuleID // id of this module
	Mempool t.ModuleID
	Net     t.ModuleID
	Crypto  t.ModuleID
}

func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:    "availability",
		Mempool: "mempool",
		Net:     "net",
		Crypto:  "crypto",
	}
}

type moduleState struct {
}

type requestMempoolBatchContext struct {
	origin *apb.RequestBatchOrigin
}

func NewModule(mc *ModuleConfig) modules.Module {
	m := dsl.NewModule(mc.Self)

	state := moduleState{}

	adsl.UponRequestBatch(m, func(origin *apb.RequestBatchOrigin) error {
		mempooldsl.RequestBatch(m, mc.Mempool, &requestMempoolBatchContext{origin})
		return nil
	})

	mempooldsl.UponNewBatch(m, func(txs [][]byte, context *requestMempoolBatchContext) error {
		submodules.Start(m)
	})

	return m
}
