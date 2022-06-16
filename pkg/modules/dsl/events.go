package dsl

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Event-specific dsl functions for emitting events (wrappers around the functions defined in pkg/events)

func SendMessage(m DslModule, destModule t.ModuleID, msg *messagepb.Message, dest []t.NodeID) {
	EmitEvent(m, events.SendMessage(destModule, msg, dest))
}

func SignRequest(m DslModule, destModule t.ModuleID, data [][]byte, origin *eventpb.SignOrigin) {
	EmitEvent(m, events.SignRequest(destModule, data, origin))
}

func VerifyNodeSigs(
	m DslModule,
	destModule t.ModuleID,
	data [][][]byte,
	signatures [][]byte,
	nodeIDs []t.NodeID,
	origin *eventpb.SigVerOrigin,
) {
	EmitEvent(m, events.VerifyNodeSigs(destModule, data, signatures, nodeIDs, origin))
}

// Event-specific dsl functions for processing events

func UponRequest(m DslModule, handler func(clientId string, reqNo uint64, data []byte, authenticator []byte) error) {
	UponEvent[eventpb.Event_Request](m, func(req *requestpb.Request) error {
		return handler(req.ClientId, req.ReqNo, req.Data, req.Authenticator)
	})
}

func UponSignResult(m DslModule, handler func(signature []byte, origin *eventpb.SignOrigin) error) {
	UponEvent[eventpb.Event_SignRequest](m, func(res *eventpb.SignResult) error {
		return handler(res.Signature, res.Origin)
	})
}

//func UponNodeSigsVerified(m DslModule, handler func(signature []byte, origin *eventpb.SigVerOrigin) error) {
//	UponEvent[eventpb.Event_NodeSigsVerified](m, func(res *eventpb.NodeSigsVerified) error {
//		return handler(res.Signature, res.Origin)
//	})
//}

func UponMessageReceived(m DslModule, handler func(from t.NodeID, msg *messagepb.Message) error) {
	UponEvent[eventpb.Event_MessageReceived](m, func(ev *eventpb.MessageReceived) error {
		return handler(t.NodeID(ev.From), ev.Msg)
	})
}
