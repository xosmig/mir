package cb

import (
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	"github.com/filecoin-project/mir/pkg/pb/cbpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func UponCBMessageReceived(m *dsl.Module, handler func(from t.NodeID, msg *cbpb.CBMessage) error) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) error {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Cb)
		if !ok {
			return nil
		}

		return handler(from, cbMsgWrapper.Cb)
	})
}

func UponStartMessageReceived(m *dsl.Module, handler func(from t.NodeID, msg *cbpb.StartMessage) error) {
	UponCBMessageReceived(m, func(from t.NodeID, msg *cbpb.CBMessage) error {
		startMsgWrapper, ok := msg.Type.(*cbpb.CBMessage_StartMessage)
		if !ok {
			return nil
		}

		return handler(from, startMsgWrapper.StartMessage)
	})
}

func UponEchoMessageReceived(m *dsl.Module, handler func(from t.NodeID, msg *cbpb.EchoMessage) error) {
	UponCBMessageReceived(m, func(from t.NodeID, msg *cbpb.CBMessage) error {
		echoMsgWrapper, ok := msg.Type.(*cbpb.CBMessage_EchoMessage)
		if !ok {
			return nil
		}

		return handler(from, echoMsgWrapper.EchoMessage)
	})
}
