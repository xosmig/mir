package cbdsl

import (
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	"github.com/filecoin-project/mir/pkg/pb/cbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Contains event wrappers specific for this module

func SignRequest(m *CBModule, destModule t.ModuleID, data [][]byte) {
	dsl.SignRequest(m, destModule, data, &eventpb.SignOrigin{
		Module: m.moduleId.Pb(),
		Type: &eventpb.SignOrigin_Empty{
			Empty: &eventpb.EmptySignOrigin{},
		},
	})
}

// Module-specific dsl handler wrappers

func UponCBMessageReceived(m *CBModule, handler func(from t.NodeID, msg *cbpb.CBMessage) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Cb)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.Cb)
	})
}

func UponStartMessageReceived(m *CBModule, handler func(from t.NodeID, msg *cbpb.StartMessage) error) {
	UponCBMessageReceived(m, func(from t.NodeID, msg *cbpb.CBMessage) error {
		startMsgWrapper, ok := msg.Type.(*cbpb.CBMessage_StartMessage)
		if !ok {
			return nil
		}

		return handler(from, startMsgWrapper.StartMessage)
	})
}

func UponEchoMessageReceived(m *CBModule, handler func(from t.NodeID, msg *cbpb.EchoMessage) error) {
	UponCBMessageReceived(m, func(from t.NodeID, msg *cbpb.CBMessage) error {
		echoMsgWrapper, ok := msg.Type.(*cbpb.CBMessage_EchoMessage)
		if !ok {
			return nil
		}

		return handler(from, echoMsgWrapper.EchoMessage)
	})
}

func VerifyNodeSignature(m *CBModule,
	destModule t.ModuleID,
	data [][]byte,
	signature []byte,
	nodeID t.NodeID,
) {
	origin := &eventpb.SigVerOrigin{
		Module: m.moduleId.Pb(),
		Type: &eventpb.SigVerOrigin_Empty{
			&eventpb.EmptySigVerOrigin{},
		},
	}
	dsl.VerifyNodeSigs(m, destModule, [][][]byte{data}, [][]byte{signature}, []t.NodeID{nodeID}, origin)
}
