package cbdsl

import (
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	"github.com/filecoin-project/mir/pkg/pb/cbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Contains event wrappers specific for this module

func SignRequest(m *cbModuleImpl, destModule t.ModuleID, data [][]byte) {
	dsl.SignRequest(m, destModule, data, &eventpb.SignOrigin{
		Module: m.moduleId.Pb(),
		Type: &eventpb.SignOrigin_Empty{
			Empty: &eventpb.EmptySignOrigin{},
		},
	})
}

// Module-specific dsl handler wrappers

func UponCBMessageReceived(m *cbModuleImpl, handler func(from t.NodeID, msg *cbpb.CBMessage) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Cb)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.Cb)
	})
}

func UponStartMessageReceived(m *cbModuleImpl, handler func(from t.NodeID, msg *cbpb.StartMessage) error) {
	UponCBMessageReceived(m, func(from t.NodeID, msg *cbpb.CBMessage) error {
		startMsgWrapper, ok := msg.Type.(*cbpb.CBMessage_StartMessage)
		if !ok {
			return nil
		}

		return handler(from, startMsgWrapper.StartMessage)
	})
}

func UponEchoMessageReceived(m *cbModuleImpl, handler func(from t.NodeID, msg *cbpb.EchoMessage) error) {
	UponCBMessageReceived(m, func(from t.NodeID, msg *cbpb.CBMessage) error {
		echoMsgWrapper, ok := msg.Type.(*cbpb.CBMessage_EchoMessage)
		if !ok {
			return nil
		}

		return handler(from, echoMsgWrapper.EchoMessage)
	})
}

func UponFinalMessageReceived(m *cbModuleImpl, handler func(from t.NodeID, msg *cbpb.FinalMessage) error) {
	UponCBMessageReceived(m, func(from t.NodeID, msg *cbpb.CBMessage) error {
		finalMsgWrapper, ok := msg.Type.(*cbpb.CBMessage_FinalMessage)
		if !ok {
			return nil
		}

		return handler(from, finalMsgWrapper.FinalMessage)
	})
}

func VerifyNodeSignature(
	m *cbModuleImpl,
	destModule t.ModuleID,
	data [][]byte,
	signature []byte,
	nodeID t.NodeID,
	origin *cbpb.CBSigVerOrigin,
) {
	dsl.VerifyNodeSigs(m, destModule, [][][]byte{data}, [][]byte{signature}, []t.NodeID{nodeID}, &eventpb.SigVerOrigin{
		Module: m.moduleId.Pb(),
		Type: &eventpb.SigVerOrigin_Cb{
			Cb: origin,
		},
	})
}

func UponCBNodeSigsVerified(m *cbModuleImpl, handler func(origin *cbpb.CBSigVerOrigin, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error) {
	dsl.UponNodeSigsVerified(m, func(origin *eventpb.SigVerOrigin, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error {
		cbOriginWrapper, ok := origin.Type.(*eventpb.SigVerOrigin_Cb)
		if !ok {
			return nil
		}

		return handler(cbOriginWrapper.Cb, nodeIds, valid, errs, allOk)
	})
}

func UponEchoSigsVerified(m *cbModuleImpl, handler func(origin *cbpb.SigVerOriginEcho, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error) {
	UponCBNodeSigsVerified(m, func(origin *cbpb.CBSigVerOrigin, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error {
		echoOriginWrapper, ok := origin.Type.(*cbpb.CBSigVerOrigin_Echo)
		if !ok {
			return nil
		}

		return handler(echoOriginWrapper.Echo, nodeIds, valid, errs, allOk)
	})
}

func UponEchoSignatureVerified(m *cbModuleImpl, handler func(origin *cbpb.SigVerOriginEcho, nodeId t.NodeID, valid bool, err error) error) {
	UponEchoSigsVerified(m, func(origin *cbpb.SigVerOriginEcho, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error {
		for i := range nodeIds {
			err := handler(origin, nodeIds[i], valid[i], errs[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}
