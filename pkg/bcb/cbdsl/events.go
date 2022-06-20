package cbdsl

import (
	"github.com/filecoin-project/mir/pkg/pb/bcbpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl handler wrappers

func UponBCBMessageReceived(m *cbModuleImpl, handler func(from t.NodeID, msg *bcbpb.BCBMessage) error) {
	dslevents.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Bcb)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.Bcb)
	})
}

func UponStartMessageReceived(m *cbModuleImpl, handler func(from t.NodeID, data []byte) error) {
	UponBCBMessageReceived(m, func(from t.NodeID, msg *bcbpb.BCBMessage) error {
		startMsgWrapper, ok := msg.Type.(*bcbpb.BCBMessage_StartMessage)
		if !ok {
			return nil
		}

		return handler(from, startMsgWrapper.StartMessage.Data)
	})
}

func UponEchoMessageReceived(m *cbModuleImpl, handler func(from t.NodeID, signature []byte) error) {
	UponBCBMessageReceived(m, func(from t.NodeID, msg *bcbpb.BCBMessage) error {
		echoMsgWrapper, ok := msg.Type.(*bcbpb.BCBMessage_EchoMessage)
		if !ok {
			return nil
		}

		return handler(from, echoMsgWrapper.EchoMessage.Signature)
	})
}

func UponFinalMessageReceived(
	m *cbModuleImpl,
	handler func(from t.NodeID, data []byte, signers []t.NodeID, signatures [][]byte) error,
) {
	UponBCBMessageReceived(m, func(from t.NodeID, msg *bcbpb.BCBMessage) error {
		finalMsgWrapper, ok := msg.Type.(*bcbpb.BCBMessage_FinalMessage)
		if !ok {
			return nil
		}

		finalMsg := finalMsgWrapper.FinalMessage

		var signers []t.NodeID
		for _, node := range finalMsg.Signers {
			signers = append(signers, t.NodeID(node))
		}

		return handler(from, finalMsg.Data, signers, finalMsg.Signatures)
	})
}

//func UponCBNodeSigsVerified(m *cbModuleImpl, handler func(origin *cbpb.CBSigVerOrigin, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error) {
//	dsl.UponNodeSigsVerified(m, func(origin *eventpb.SigVerOrigin, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error {
//		cbOriginWrapper, ok := origin.Type.(*eventpb.SigVerOrigin_Cb)
//		if !ok {
//			return nil
//		}
//
//		return handler(cbOriginWrapper.Cb, nodeIds, valid, errs, allOk)
//	})
//}

//func UponEchoSigsVerified(m *cbModuleImpl, handler func(origin *cbpb.SigVerOriginEcho, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error) {
//	UponCBNodeSigsVerified(m, func(origin *cbpb.CBSigVerOrigin, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error {
//		echoOriginWrapper, ok := origin.Type.(*cbpb.CBSigVerOrigin_Echo)
//		if !ok {
//			return nil
//		}
//
//		return handler(echoOriginWrapper.Echo, nodeIds, valid, errs, allOk)
//	})
//}

//func UponEchoSignatureVerified(m *cbModuleImpl, handler func(origin *cbpb.SigVerOriginEcho, nodeId t.NodeID, valid bool, err error) error) {
//	UponEchoSigsVerified(m, func(origin *cbpb.SigVerOriginEcho, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error {
//		for i := range nodeIds {
//			err := handler(origin, nodeIds[i], valid[i], errs[i])
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//}
