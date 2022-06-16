package dsl

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Event-specific dsl functions (wrappers around the functions defined in pkg/events)

func SendMessage(m *Module, destModule t.ModuleID, msg *messagepb.Message, dest []t.NodeID) {
	EmitEvent(m, events.SendMessage(destModule, msg, dest))
}

func UponMessageReceived(m *Module, handler func(from t.NodeID, msg *messagepb.Message) error) {
	UponEvent[eventpb.Event_MessageReceived](m, func(ev *eventpb.MessageReceived) error {
		return handler(t.NodeID(ev.From), ev.Msg)
	})
}

func SignRequest(m *Module, destModule t.ModuleID, data [][]byte, origin *eventpb.SignOrigin) {
	EmitEvent(m, events.SignRequest(destModule, data, origin))
}

func VerifyNodeSigs(
	m *Module,
	destModule t.ModuleID,
	data [][][]byte,
	signatures [][]byte,
	nodeIDs []t.NodeID,
	origin *eventpb.SigVerOrigin,
) {
	EmitEvent(m, events.VerifyNodeSigs(destModule, data, signatures, nodeIDs, origin))
}
