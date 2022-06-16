package dsl

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/pkg/errors"
)

// Event-specific dsl functions for emitting events (wrappers around the functions defined in pkg/events)

func SendMessage(m DslModule, destModule t.ModuleID, msg *messagepb.Message, dest []t.NodeID) {
	EmitEvent(m, events.SendMessage(destModule, msg, dest))
}

func SignRequest(m DslModule, destModule t.ModuleID, data [][]byte, origin *eventpb.SignOrigin) {
	EmitEvent(m, events.SignRequest(destModule, data, origin))
}

// VerifyNodeSigs emits a signature verification event for a batch of signatures.
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

func UponRequest(m DslModule, handler func(clientId t.ClientID, reqNo uint64, data []byte, authenticator []byte) error) {
	UponEvent[eventpb.Event_Request](m, func(req *requestpb.Request) error {
		return handler(t.ClientID(req.ClientId), req.ReqNo, req.Data, req.Authenticator)
	})
}

func UponSignResult(m DslModule, handler func(signature []byte, origin *eventpb.SignOrigin) error) {
	UponEvent[eventpb.Event_SignRequest](m, func(res *eventpb.SignResult) error {
		return handler(res.Signature, res.Origin)
	})
}

func UponNodeSigsVerified(m DslModule, handler func(origin *eventpb.SigVerOrigin, nodeIds []t.NodeID, valid []bool, errs []error, allOk bool) error) {
	UponEvent[eventpb.Event_NodeSigsVerified](m, func(res *eventpb.NodeSigsVerified) error {
		var nodeIds []t.NodeID
		for _, id := range res.NodeIds {
			nodeIds = append(nodeIds, t.NodeID(id))
		}

		var errs []error
		for _, err := range res.Errors {
			errs = append(errs, errors.New(err))
		}

		return handler(res.Origin, nodeIds, res.Valid, errs, res.AllOk)
	})
}

func UponMessageReceived(m DslModule, handler func(from t.NodeID, msg *messagepb.Message) error) {
	UponEvent[eventpb.Event_MessageReceived](m, func(ev *eventpb.MessageReceived) error {
		return handler(t.NodeID(ev.From), ev.Msg)
	})
}
