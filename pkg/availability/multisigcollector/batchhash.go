package multisigcollector

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
)

func processEventsForComputingBatchHash(m dsl.Module, mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) {
	// When the hashes of the received transactions are computed, compute the hash of the batch.
	dsl.UponHashResult(m, func(hashes [][]byte, context *computeHashOfReceivedTxsContext) error {
		dsl.OneHashRequest(m, mc.Hasher, hashes, &computeHashOfReceivedBatchContext{context.sourceID, context.reqID})
		return nil
	})

	// When the hash of the batch is computed, persistently store the batch and generate a signature.
	dsl.UponOneHashResult(m, func(batchHash []byte, context *computeHashOfReceivedBatchContext) error {
		sigMsg := sigMessage(params.InstanceUID, context.sourceID, context.reqID, batchHash)
		dsl.SignRequest(m, mc.Crypto, sigMsg, &signReceivedBatchContext{context.sourceID, context.reqID})
		return nil
	})

	processEventsForComputingBatchHash(m, mc, params, nodeID)
	processEventsForCreatingCertificates(m, mc, params, nodeID)
	processEventsForRetrievingBatches(m, mc, params, nodeID)
}
