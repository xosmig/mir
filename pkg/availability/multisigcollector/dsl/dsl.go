package Mscdsl

import (
	"fmt"
	adsl "github.com/filecoin-project/mir/pkg/availability/dsl"
	msc "github.com/filecoin-project/mir/pkg/availability/multisigcollector"
	cs "github.com/filecoin-project/mir/pkg/contextstore"
	"github.com/filecoin-project/mir/pkg/dsl"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponMscMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *mscpb.Message) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_MultisigCollector)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.MultisigCollector)
	})
}

func UponRequestSigMessageReceived(m dsl.Module, handler func(from t.NodeID, txs [][]byte, id cs.ItemID) error) {
	UponMscMessageReceived(m, func(from t.NodeID, msg *mscpb.Message) error {
		requestSigMsgWrapper, ok := msg.Type.(*mscpb.Message_RequestSig)
		if !ok {
			return nil
		}
		requestSigMsg := requestSigMsgWrapper.RequestSig

		return handler(from, requestSigMsg.Txs, cs.ItemID(requestSigMsg.CsItemId))
	})
}

func UponSigMessageReceived(m dsl.Module, handler func(from t.NodeID, reqID msc.RequestID, signature []byte) error) {
	UponMscMessageReceived(m, func(from t.NodeID, msg *mscpb.Message) error {
		sigMsgWrapper, ok := msg.Type.(*mscpb.Message_Sig)
		if !ok {
			return nil
		}
		sigMsg := sigMsgWrapper.Sig

		return handler(from, msc.RequestID(sigMsg.ReqId), sigMsg.Signature)
	})
}

func UponRequestBatchMessageReceived(m dsl.Module, handler func(from t.NodeID, sourceID t.NodeID, certReqID msc.RequestID, csItemID cs.ItemID) error) {
	UponMscMessageReceived(m, func(from t.NodeID, msg *mscpb.Message) error {
		requestBatchMsgWrapper, ok := msg.Type.(*mscpb.Message_RequestBatch)
		if !ok {
			return nil
		}
		requestBatchMsg := requestBatchMsgWrapper.RequestBatch

		return handler(from, t.NodeID(requestBatchMsg.SourceId), msc.RequestID(requestBatchMsg.CertReqId), cs.ItemID(requestBatchMsg.CsItemId))
	})
}

func UponProvideBatchMessageReceived(m dsl.Module, handler func(from t.NodeID, txs [][]byte, csItemID cs.ItemID) error) {
	UponMscMessageReceived(m, func(from t.NodeID, msg *mscpb.Message) error {
		provideBatchMsgWrapper, ok := msg.Type.(*mscpb.Message_ProvideBatch)
		if !ok {
			return nil
		}
		provideBatchMsg := provideBatchMsgWrapper.ProvideBatch

		return handler(from, provideBatchMsg.Txs, cs.ItemID(provideBatchMsg.CsItemId))
	})
}

func UponRequestTransactions(m dsl.Module, handler func(cert *mscpb.Cert, origin *apb.RequestTransactionsOrigin) error) {
	adsl.UponRequestTransactions(m, func(cert *apb.Cert, origin *apb.RequestTransactionsOrigin) error {
		mscCertWrapper, ok := cert.Type.(*apb.Cert_Msc)
		if !ok {
			return fmt.Errorf("unexpected certificate type. Expected: %T, got: %T", mscCertWrapper, cert.Type)
		}

		return handler(mscCertWrapper.Msc, origin)
	})
}
