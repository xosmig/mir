package multisigcollector

import (
	"github.com/filecoin-project/mir/pkg/availability/cft-storage/internal/common"
	"github.com/filecoin-project/mir/pkg/availability/cft-storage/internal/parts/batchreconstruction"
	"github.com/filecoin-project/mir/pkg/availability/cft-storage/internal/parts/certcreation"
	"github.com/filecoin-project/mir/pkg/availability/cft-storage/internal/parts/certverification"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig = common.ModuleConfig

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:    "availability",
		Mempool: "mempool",
		Net:     "net",
		Crypto:  "crypto",
	}
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams = common.ModuleParams

// NewModule creates a new instance of the multisig collector module.
// Multisig collector is the simplest implementation of the availability layer.
// Whenever an availability certificate is requested, it pulls a batch from the mempool module,
// sends it to all replicas and collects a quorum (i.e., more than (N+F)/2) of signatures confirming that
// other nodes have persistently stored the batch.
func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	commonState := &common.State{
		BatchStore:       make(map[t.BatchID][]t.TxID),
		TransactionStore: make(map[t.TxID][]byte),
	}

	certcreation.IncludeCreatingCertificates(m, mc, params, nodeID, commonState)
	certverification.IncludeVerificationOfCertificates(m, mc, params, nodeID, commonState)
	batchreconstruction.IncludeBatchReconstruction(m, mc, params, nodeID, commonState)

	return m
}
