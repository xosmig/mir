package multisigcollector

import (
	cs "github.com/filecoin-project/mir/pkg/contextstore"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Message(moduleID t.ModuleID, msg *mscpb.Message) *messagepb.Message {
	return &messagepb.Message{
		DestModule: moduleID.Pb(),
		Type: &messagepb.Message_MultisigCollector{
			MultisigCollector: msg,
		},
	}
}

func RequestSigMessage(moduleID t.ModuleID, txs [][]byte, csItemID cs.ItemID) *messagepb.Message {
	return Message(moduleID, &mscpb.Message{
		Type: &mscpb.Message_RequestSig{
			RequestSig: &mscpb.RequestSigMessage{
				Txs:      txs,
				CsItemId: csItemID.Pb(),
			},
		},
	})
}

func EchoMessage(moduleID t.ModuleID, reqID RequestID, signature []byte) *messagepb.Message {
	return Message(moduleID, &mscpb.Message{
		Type: &mscpb.Message_Sig{
			Sig: &mscpb.SigMessage{
				ReqId:     reqID.Pb(),
				Signature: signature,
			},
		},
	})
}

func RequestBatchMessage(moduleID t.ModuleID, sourceID t.NodeID, certReqID RequestID, csItemID cs.ItemID) *messagepb.Message {
	return Message(moduleID, &mscpb.Message{
		Type: &mscpb.Message_RequestBatch{
			RequestBatch: &mscpb.RequestBatchMessage{
				SourceId:  sourceID.Pb(),
				CertReqId: certReqID.Pb(),
				CsItemId:  csItemID.Pb(),
			},
		},
	})
}

func ProvideBatchMessage(moduleID t.ModuleID, txs [][]byte, csItemID cs.ItemID) *messagepb.Message {
	return Message(moduleID, &mscpb.Message{
		Type: &mscpb.Message_ProvideBatch{
			ProvideBatch: &mscpb.ProvideBatchMessage{
				Txs:      txs,
				CsItemId: csItemID.Pb(),
			},
		},
	})
}

func Cert(sourceID t.NodeID, reqID RequestID, batchHash []byte, signers []t.NodeID, signatures [][]byte) *apb.Cert {
	return &apb.Cert{
		Type: &apb.Cert_Msc{
			Msc: &mscpb.Cert{
				SourceId:   sourceID.Pb(),
				ReqId:      reqID.Pb(),
				BatchHash:  batchHash,
				Signers:    t.NodeIDSlicePb(signers),
				Signatures: signatures,
			},
		},
	}
}
