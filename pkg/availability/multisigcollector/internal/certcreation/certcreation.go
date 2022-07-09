package certcreation

import (
	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	msc "github.com/filecoin-project/mir/pkg/availability/multisigcollector"
	mscdsl "github.com/filecoin-project/mir/pkg/availability/multisigcollector/dsl"
	"github.com/filecoin-project/mir/pkg/availability/multisigcollector/internal/common"
	"github.com/filecoin-project/mir/pkg/dsl"
	mempooldsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

type localState struct {
	*common.State
	nextReqID   msc.RequestID
	sourceState map[msc.RequestID]*sourceState
}

// sourceState represents the state of the broadcaster (leader) in an instance of consistent broadcast.
// The source can dispose of the state of the request as soon as the request is completed.
type sourceState struct {
	reqOrigin *apb.RequestCertOrigin
	batchHash []byte

	receivedSig map[t.NodeID]bool
	sigs        map[t.NodeID][]byte
}

func ProcessEventsForCreatingCertificates(
	m dsl.Module,
	mc *msc.ModuleConfig,
	params *msc.ModuleParams,
	commonState *common.State,
) {
	state := localState{
		State:       commonState,
		nextReqID:   0,
		sourceState: make(map[msc.RequestID]*sourceState),
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// creating and storing the batch                                                                                 //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// When a batch is requested by the consensus layer, initialize an instance of the broadcast protocol and ask
	// mempool for a batch.
	adsl.UponRequestCert(m, func(origin *apb.RequestCertOrigin) error {
		reqID := state.nextReqID
		state.nextReqID++

		state.sourceState[reqID] = &sourceState{
			reqOrigin:   origin,
			receivedSig: make(map[t.NodeID]bool),
			sigs:        make(map[t.NodeID][]byte),
		}

		mempooldsl.RequestBatch(m, mc.Mempool, &requestBatchFromMempoolContext{reqID})
		return nil
	})

	// When the mempool provides a batch, compute its hash.
	mempooldsl.UponNewBatch(m, func(txIDs [][]byte, txs [][]byte, context *requestBatchFromMempoolContext) error {
		dsl.OneHashRequest(m, mc.Hasher, txIDs, &computeHashOfOwnBatchContext{context.reqID, txs})
		return nil
	})

	// When the hash of the batch is computed, request signatures for the batch from all nodes.
	dsl.UponOneHashResult(m, func(batchHash []byte, context *computeHashOfOwnBatchContext) error {
		state.sourceState[context.reqID].batchHash = batchHash
		// TODO: add persistent storage for crash-recovery.
		dsl.SendMessage(m, mc.Net, msc.RequestSigMessage(mc.Self, context.reqID, context.txs), params.AllNodes)
		return nil
	})

	// When receive a request for a signature, store the received transactions and compute their hashes.
	mscdsl.UponRequestSigMessageReceived(m, func(from t.NodeID, reqID msc.RequestID, txs [][]byte) error {
		if _, ok := state.replicaState[from][reqID]; ok {
			// Already received a request with the same source and reqID.
			return nil
		}
		// TODO: replicaState should be persisted for crash-recovery.
		state.replicaState[from][reqID] = &replicaState{
			sentSig: false,
			txs:     txs,
		}

		txsMsgs := sliceutil.Transform(txs, func(tx []byte) [][]byte { return [][]byte{tx} })
		dsl.HashRequest(m, mc.Hasher, txsMsgs, &computeHashOfReceivedTxsContext{from, reqID})
		return nil
	})

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

	// When the hashes of the received transactions are computed, compute the hash of the batch.
	dsl.UponHashResult(m, func(hashes [][]byte, context *computeHashOfReceivedTxsOnSigRequestContext) error {
		dsl.OneHashRequest(m, mc.Hasher, hashes, &computeHashOfReceivedBatchOnSigRequestContext{context.sourceID, context.reqID})
		return nil
	})

	// When the hash of the batch is computed, persistently store the batch and generate a signature.
	dsl.UponOneHashResult(m, func(batchHash []byte, context *computeHashOfReceivedBatchOnSigRequestContext) error {
		sigMsg := sigMessage(params.InstanceUID, context.sourceID, context.reqID, batchHash)
		dsl.SignRequest(m, mc.Crypto, sigMsg, &signReceivedBatchOnSigRequestContext{context.sourceID, context.reqID})
		return nil
	})

	// When a signature is generated, send it to the process that sent the request.
	dsl.UponSignResult(m, func(signature []byte, context *signReceivedBatchOnSigRequestContext) error {
		requestState := state.replicaState[context.sourceID][context.reqID]
		if !requestState.sentSig {
			requestState.sentSig = true
			dsl.SendMessage(m, mc.Net, EchoMessage(mc.Self, context.reqID, signature), []t.NodeID{context.sourceID})
		}
		return nil
	})

	// When receive a signature, verify its correctness.
	mscdsl.UponSigMessageReceived(m, func(from t.NodeID, reqID msc.RequestID, signature []byte) error {
		requestState, ok := state.sourceState[reqID]
		if !ok {
			// Ignore a message with an invalid or outdated request id.
			return nil
		}

		if !requestState.receivedSig[from] {
			requestState.receivedSig[from] = true
			sigMsg := sigMessage(params.InstanceUID, nodeID, reqID, requestState.batchHash)
			dsl.VerifyOneNodeSig(m, mc.Crypto, sigMsg, signature, from, &verifySigContext{reqID, signature})
		}
		return nil
	})

	// When a signature is verified, store it in memory.
	dsl.UponOneNodeSigVerified(m, func(nodeID t.NodeID, err error, context *verifySigContext) error {
		if err != nil {
			// Ignore invalid signature.
			return nil
		}
		requestState, ok := state.sourceState[context.reqID]
		if !ok {
			// The request has already been completed.
			return nil
		}

		requestState.sigs[nodeID] = context.signature
		return nil
	})

	// When a quorum (more than (N+F)/2) of signatures are collected, create and output a certificate.
	dsl.UponCondition(m, func() error {
		// Iterate over active leader instances.
		// Note that, most of the time, there will be at most one active instance at a time.
		for reqID, requestState := range state.sourceState {
			if len(requestState.sigs) > (params.N()+params.F())/2 {
				certNodes, certSigs := maputil.GetKeysAndValues(requestState.sigs)

				requestingModule := t.ModuleID(requestState.reqOrigin.Module)
				cert := Cert(nodeID, reqID, requestState.batchHash, certNodes, certSigs)
				adsl.NewCert(m, requestingModule, cert, requestState.reqOrigin)

				// Dispose of the state associated with this instance.
				// Note that the replicas cannot dispose of their state.
				delete(state.sourceState, reqID)
			}
		}
		return nil
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Context data structures                                                                                            //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type requestBatchFromMempoolContext struct {
	reqID msc.RequestID
}

type computeHashOfOwnBatchContext struct {
	reqID msc.RequestID
	txs   [][]byte
}

type computeHashOfReceivedTxsContext struct {
	sourceID t.NodeID
	reqID    msc.RequestID
}

type computeHashOfReceivedBatchContext struct {
	sourceID t.NodeID
	reqID    msc.RequestID
}

type verifySigContext struct {
	reqID     msc.RequestID
	signature []byte
}
