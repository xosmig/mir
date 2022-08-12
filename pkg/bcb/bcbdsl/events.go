package bcbdsl

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/bcbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func Request(m dsl.Module, dest t.ModuleID, data []byte) {
	dsl.EmitEvent(m, &eventpb.Event{
		DestModule: dest.Pb(),

		Type: &eventpb.Event_Bcb{
			Bcb: &bcbpb.Event{
				Type: &bcbpb.Event_Request{
					Request: &bcbpb.BroadcastRequest{
						Data: data,
					},
				},
			},
		},
	})
}

func Deliver(m dsl.Module, dest t.ModuleID, data []byte) {
	dsl.EmitEvent(m, &eventpb.Event{
		DestModule: dest.Pb(),

		Type: &eventpb.Event_Bcb{
			Bcb: &bcbpb.Event{
				Type: &bcbpb.Event_Deliver{
					Deliver: &bcbpb.Deliver{
						Data: data,
					},
				},
			},
		},
	})
}

// Module-specific dsl functions for processing events.

func UponEvent[EvWrapper bcbpb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) dsl.Result) {
	dsl.UponEvent[*eventpb.Event_Bcb](m, func(ev *bcbpb.Event) dsl.Result {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return dsl.NoMatch()
		}
		return handler(evWrapper.Unwrap())
	})
}

func UponBroadcastRequest(m dsl.Module, handler func(data []byte) dsl.Result) {
	UponEvent[*bcbpb.Event_Request](m, func(ev *bcbpb.BroadcastRequest) dsl.Result {
		return handler(ev.Data)
	})
}

func UponDeliver(m dsl.Module, handler func(data []byte) dsl.Result) {
	UponEvent[*bcbpb.Event_Deliver](m, func(ev *bcbpb.Deliver) dsl.Result {
		return handler(ev.Data)
	})
}

func UponBCBMessageReceived(m dsl.Module, handler func(from t.NodeID, msg *bcbpb.Message) dsl.Result) {
	dsl.UponMessageReceived(m, func(from t.NodeID, msg *messagepb.Message) dsl.Result {
		cbMsgWrapper, ok := msg.Type.(*messagepb.Message_Bcb)
		if !ok {
			return dsl.NoMatch()
		}

		return handler(from, cbMsgWrapper.Bcb)
	})
}

func UponStartMessageReceived(m dsl.Module, handler func(from t.NodeID, data []byte) dsl.Result) {
	UponBCBMessageReceived(m, func(from t.NodeID, msg *bcbpb.Message) dsl.Result {
		startMsgWrapper, ok := msg.Type.(*bcbpb.Message_StartMessage)
		if !ok {
			return dsl.NoMatch()
		}

		return handler(from, startMsgWrapper.StartMessage.Data)
	})
}

func UponEchoMessageReceived(m dsl.Module, handler func(from t.NodeID, signature []byte) dsl.Result) {
	UponBCBMessageReceived(m, func(from t.NodeID, msg *bcbpb.Message) dsl.Result {
		echoMsgWrapper, ok := msg.Type.(*bcbpb.Message_EchoMessage)
		if !ok {
			return dsl.NoMatch()
		}

		return handler(from, echoMsgWrapper.EchoMessage.Signature)
	})
}

func UponFinalMessageReceived(
	m dsl.Module,
	handler func(from t.NodeID, data []byte, signers []t.NodeID, signatures [][]byte) dsl.Result,
) {
	UponBCBMessageReceived(m, func(from t.NodeID, msg *bcbpb.Message) dsl.Result {
		finalMsgWrapper, ok := msg.Type.(*bcbpb.Message_FinalMessage)
		if !ok {
			return dsl.NoMatch()
		}

		finalMsg := finalMsgWrapper.FinalMessage
		return handler(from, finalMsg.Data, t.NodeIDSlice(finalMsg.Signers), finalMsg.Signatures)
	})
}
