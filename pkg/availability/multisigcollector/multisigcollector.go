package multisigcollector

import (
	"encoding/binary"
	"fmt"
	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	mscdsl "github.com/filecoin-project/mir/pkg/availability/multisigcollector/dsl"
	cs "github.com/filecoin-project/mir/pkg/contextstore"
	"github.com/filecoin-project/mir/pkg/dsl"
	mempooldsl "github.com/filecoin-project/mir/pkg/mempool/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self    t.ModuleID // id of this module
	Mempool t.ModuleID
	Net     t.ModuleID
	Crypto  t.ModuleID
	Hasher  t.ModuleID
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:    "availability",
		Mempool: "mempool",
		Net:     "net",
		Crypto:  "crypto",
		Hasher:  "hasher",
	}
}

// InstanceUID is used to uniquely identify an instance of multisig collector.
// It is used to prevent cross-instance signature replay attack and should be unique across all executions.
type InstanceUID []byte

// Bytes returns the binary representation of the InstanceUID.
func (uid InstanceUID) Bytes() []byte {
	return uid
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	InstanceUID []byte     // unique identifier for this instance of BCB, used to prevent cross-instance replay attacks
	AllNodes    []t.NodeID // the list of participating nodes
}

// N is the total number of replicas.
func (params *ModuleParams) N() int {
	return len(params.AllNodes)
}

// F is the maximum number of replicas that can be tolerated.
func (params *ModuleParams) F() int {
	return (params.N() - 1) / 3
}

// RequestID is used to uniquely identify multisig collector requests.
type RequestID uint64

// Pb returns the protobuf representation of a RequestID.
func (id RequestID) Pb() uint64 {
	return uint64(id)
}

// Bytes returns the binary representation of the RequestID.
func (id RequestID) Bytes() []byte {
	var res []byte
	binary.LittleEndian.PutUint64(res, uint64(id))
	return res
}

// requestSourceState represents the state of the broadcaster (leader) in an instance of consistent broadcast.
// The source can dispose of the state of the request as soon as the request is completed.
type requestSourceState struct {
	reqOrigin *apb.RequestCertOrigin
	batchHash []byte

	receivedSig map[t.NodeID]bool
	sigs        map[t.NodeID][]byte
}

// requestReplicaState represents the state of a non-broadcaster in an instance of consistent broadcast.
// A replica must never dispose of its state because in order to not sign a contradicting message and to be able to
// recover the transactions.
type requestReplicaState struct {
	sentSig bool
	txs     [][]byte
}

// moduleState represents the state of the multisig collector module.
type moduleState struct {
	nextReqID    RequestID
	sourceState  map[RequestID]*requestSourceState
	replicaState map[t.NodeID]map[RequestID]*requestReplicaState

	requestBatchMsgContextStore cs.ContextStore[*requestBatchMsgContext]
}

// NewModule creates a new instance of the multisig collector module.
func NewModule(mc *ModuleConfig, params *ModuleParams, nodeID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	// initialize the state of the module
	state := moduleState{
		nextReqID:    0,
		sourceState:  make(map[RequestID]*requestSourceState),
		replicaState: make(map[t.NodeID]map[RequestID]*requestReplicaState),
	}

	for _, id := range params.AllNodes {
		state.replicaState[id] = make(map[RequestID]*requestReplicaState)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// creating and storing the batch                                                                                 //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// When a batch is requested by the consensus layer, initialize an instance of the broadcast protocol and ask
	// mempool for a batch.
	adsl.UponRequestCert(m, func(origin *apb.RequestCertOrigin) error {
		reqID := state.nextReqID
		state.nextReqID++

		state.sourceState[reqID] = &requestSourceState{
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
		dsl.SendMessage(m, mc.Net, RequestSigMessage(mc.Self, context.reqID, context.txs), params.AllNodes)
		return nil
	})

	// When receive a request for a signature, store the received transactions and compute their hashes.
	mscdsl.UponRequestSigMessageReceived(m, func(from t.NodeID, reqID RequestID, txs [][]byte) error {
		if _, ok := state.replicaState[from][reqID]; ok {
			// Already received a request with the same source and reqID.
			return nil
		}
		// TODO: replicaState should be persisted for crash-recovery.
		state.replicaState[from][reqID] = &requestReplicaState{
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

	// When a signature is generated, send it to the process that sent the request.
	dsl.UponSignResult(m, func(signature []byte, context *signReceivedBatchContext) error {
		requestState := state.replicaState[context.sourceID][context.reqID]
		if !requestState.sentSig {
			requestState.sentSig = true
			dsl.SendMessage(m, mc.Net, EchoMessage(mc.Self, context.reqID, signature), []t.NodeID{context.sourceID})
		}
		return nil
	})

	// When receive a signature, verify its correctness.
	mscdsl.UponSigMessageReceived(m, func(from t.NodeID, reqID RequestID, signature []byte) error {
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

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// verifying correctness of a certificate                                                                         //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// When receive a request to verify a certificate, check that it is structurally correct and verify the signatures.
	adsl.UponVerifyCert(m, func(cert *apb.Cert, origin *apb.VerifyCertOrigin) error {
		mscCert, err := verifyCertificateStructure(params, cert)
		if err != nil {
			adsl.CertVerified(m, t.ModuleID(origin.Module), false, origin)
			return nil
		}

		sigMsg := sigMessage(params.InstanceUID, t.NodeID(mscCert.SourceId), RequestID(mscCert.ReqId), mscCert.BatchHash)
		dsl.VerifyNodeSigs(m, mc.Crypto,
			/*data*/ sliceutil.Repeat(sigMsg, len(mscCert.Signers)),
			/*signatures*/ mscCert.Signatures,
			/*nodeIDs*/ t.NodeIDSlice(mscCert.Signers),
			/*context*/ &verifySigsInCertContext{origin},
		)
		return nil
	})

	// When the signatures in a certificate are verified, output the result of certificate verification.
	dsl.UponNodeSigsVerified(m, func(nodeIDs []t.NodeID, errs []error, allOK bool, context *verifySigsInCertContext) error {
		adsl.CertVerified(m, t.ModuleID(context.origin.Module), allOK, context.origin)
		return nil
	})

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// recovering the batch by its certificate                                                                        //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// When receive a request for transactions, first check the local persistent storage and then ask other nodes.
	mscdsl.UponRequestTransactions(m, func(cert *mscpb.Cert, origin *apb.RequestTransactionsOrigin) error {
		requestState, ok := state.replicaState[t.NodeID(cert.SourceId)][RequestID(cert.ReqId)]
		if ok {
			adsl.ProvideTransactions(m, t.ModuleID(origin.Module), requestState.txs, origin)
			return nil
		}

		csItemID := state.requestBatchMsgContextStore.Store(&requestBatchMsgContext{origin: origin})
		dsl.SendMessage(m, mc.Net,
			RequestBatchMessage(mc.Self, t.NodeID(cert.SourceId), RequestID(cert.ReqId), csItemID),
			t.NodeIDSlice(cert.Signers))
		return nil
	})

	// When receive a request for batch from another node, send all transactions in response.
	mscdsl.UponRequestBatchMessageReceived(m, func(from t.NodeID, sourceID t.NodeID, certReqID RequestID, csItemID cs.ItemID) error {
		requestState, ok := state.replicaState[sourceID][certReqID]
		if !ok {
			// Ignore invalid request.
			return nil
		}

		dsl.SendMessage(m, mc.Net, ProvideBatchMessage(mc.Self, requestState.txs, csItemID), []t.NodeID{from})
		return nil
	})

	// When receive a requested batch, compute the hashes of the received transactions.
	mscdsl.UponProvideBatchMessageReceived(m, func(from t.NodeID, txs [][]byte, csItemID cs.ItemID) error {
		context, ok := state.requestBatchMsgContextStore.Recover(csItemID)
		if !ok {
			// Ignore a message with an invalid or outdated csItemID.
			return nil
		}

		txsMsgs := sliceutil.Transform(txs, func(tx []byte) [][]byte { return [][]byte{tx} })
		dsl.HashRequest(m, mc.Hasher, txsMsgs, &computeHash)
	})

	return m
}

func verifyCertificateStructure(params *ModuleParams, cert *apb.Cert) (*mscpb.Cert, error) {
	// Check that the certificate is present.
	if cert == nil || cert.Type == nil {
		return nil, fmt.Errorf("the certificate is nil")
	}

	// Check that the certificate is of the right type.
	mscCertWrapper, ok := cert.Type.(*apb.Cert_Msc)
	if !ok {
		return nil, fmt.Errorf("unexpected certificate type")
	}
	mscCert := mscCertWrapper.Msc

	// Check that the certificate contains a sufficient number of signatures.
	if len(mscCert.Signers) <= (params.N()+params.F())/2 {
		return nil, fmt.Errorf("insuficient number of signatures")
	}

	if len(mscCert.Signers) != len(mscCert.Signatures) {
		return nil, fmt.Errorf("the number of signatures does not correspond to the number of signers")
	}

	// Check that the identities of the signing nodes are not repeated.
	alreadySeen := make(map[t.NodeID]struct{})
	for _, idRaw := range mscCert.Signers {
		id := t.NodeID(idRaw)
		if _, ok := alreadySeen[id]; ok {
			return nil, fmt.Errorf("some node ids in the certificate are repeated multiple times")
		}
		alreadySeen[id] = struct{}{}
	}

	// Check that the identities of the source node and the signing nodes are valid.
	allNodes := make(map[t.NodeID]struct{})
	for _, id := range params.AllNodes {
		allNodes[id] = struct{}{}
	}

	if _, ok := allNodes[t.NodeID(mscCert.SourceId)]; !ok {
		return nil, fmt.Errorf("unknown source node id: %v", t.NodeID(mscCert.SourceId))
	}

	for _, idRaw := range mscCert.Signers {
		if _, ok := allNodes[t.NodeID(idRaw)]; !ok {
			return nil, fmt.Errorf("unknown node id: %v", t.NodeID(idRaw))
		}
	}

	return mscCert, nil
}

func sigMessage(instanceUID InstanceUID, sourceID t.NodeID, reqID RequestID, batchHash []byte) [][]byte {
	return [][]byte{instanceUID.Bytes(), sourceID.Bytes(), reqID.Bytes(), batchHash}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Context data structures                                                                                            //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type requestBatchFromMempoolContext struct {
	reqID RequestID
}

type computeHashOfOwnBatchContext struct {
	reqID RequestID
	txs   [][]byte
}

type computeHashOfReceivedTxsContext struct {
	sourceID t.NodeID
	reqID    RequestID
}

type computeHashOfReceivedBatchContext struct {
	sourceID t.NodeID
	reqID    RequestID
}

type signReceivedBatchContext struct {
	sourceID t.NodeID
	reqID    RequestID
}

type verifySigContext struct {
	reqID     RequestID
	signature []byte
}

type verifySigsInCertContext struct {
	origin *apb.VerifyCertOrigin
}

type requestBatchMsgContext struct {
	origin *apb.RequestTransactionsOrigin
}
