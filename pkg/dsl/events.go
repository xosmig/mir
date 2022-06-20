package dsl

import (
	cs "github.com/filecoin-project/mir/pkg/contextstore"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/pkg/errors"
)

// Event-specific dsl functions for emitting events (wrappers around the functions defined in pkg/events)

func SendMessage(m Module, destModule t.ModuleID, msg *messagepb.Message, dest []t.NodeID) {
	EmitEvent(m, events.SendMessage(destModule, msg, dest))
}

func SignRequest[C any](m Module, destModule t.ModuleID, data [][]byte, contextHandle ContextHandle[C], context C) {
	itemID := contextHandle.StoreContext(context)

	origin := &eventpb.SignOrigin{
		Module: m.GetModuleID().Pb(),
		Type: &eventpb.SignOrigin_Dsl{
			Dsl: &eventpb.DslOrigin{
				ContextHandleId: contextHandle.GetID().Pb(),
				ItemID:          itemID.Pb(),
			},
		},
	}

	EmitEvent(m, events.SignRequest(destModule, data, origin))
}

// VerifyNodeSigs emits a signature verification event for a batch of signatures.
func VerifyNodeSigs[C any](
	m Module,
	destModule t.ModuleID,
	data [][][]byte,
	signatures [][]byte,
	nodeIDs []t.NodeID,
	contextHandle ContextHandle[C],
	context C,
) {
	itemID := contextHandle.StoreContext(context)

	origin := &eventpb.SigVerOrigin{
		Module: m.GetModuleID().Pb(),
		Type: &eventpb.SigVerOrigin_Dsl{
			Dsl: &eventpb.DslOrigin{
				ContextHandleId: contextHandle.GetID().Pb(),
				ItemID:          itemID.Pb(),
			},
		},
	}

	EmitEvent(m, events.VerifyNodeSigs(destModule, data, signatures, nodeIDs, origin))
}

// Event-specific dsl functions for processing events

func UponRequest(m Module, handler func(clientId t.ClientID, reqNo uint64, data []byte, authenticator []byte) error) {
	UponEvent[eventpb.Event_Request](m, func(req *requestpb.Request) error {
		return handler(t.ClientID(req.ClientId), req.ReqNo, req.Data, req.Authenticator)
	})
}

func UponMessageReceived(m Module, handler func(from t.NodeID, msg *messagepb.Message) error) {
	UponEvent[eventpb.Event_MessageReceived](m, func(ev *eventpb.MessageReceived) error {
		return handler(t.NodeID(ev.From), ev.Msg)
	})
}

func UponSignResult[C any](m Module, contextHandle ContextHandle[C], handler func(signature []byte, context C) error) {
	UponEvent[eventpb.Event_SignRequest](m, func(res *eventpb.SignResult) error {
		dslOriginWrapper, ok := res.Origin.Type.(*eventpb.SignOrigin_Dsl)
		if !ok {
			return nil
		}
		dslOrigin := dslOriginWrapper.Dsl

		if ContextHandleID(dslOrigin.ContextHandleId) != contextHandle.GetID() {
			return nil
		}

		context := contextHandle.RecoverContextAndDeferCleanup(cs.ItemID(dslOrigin.ItemID))
		return handler(res.Signature, context)
	})
}

func UponNodeSigsVerified[C any](
	m Module,
	contextHandle ContextHandle[C],
	handler func(nodeIds []t.NodeID, valid []bool, errs []error, allOk bool, context C) error,
) {
	UponEvent[eventpb.Event_NodeSigsVerified](m, func(res *eventpb.NodeSigsVerified) error {
		dslOriginWrapper, ok := res.Origin.Type.(*eventpb.SigVerOrigin_Dsl)
		if !ok {
			return nil
		}
		dslOrigin := dslOriginWrapper.Dsl

		if ContextHandleID(dslOrigin.ContextHandleId) != contextHandle.GetID() {
			return nil
		}

		var nodeIds []t.NodeID
		for _, id := range res.NodeIds {
			nodeIds = append(nodeIds, t.NodeID(id))
		}

		var errs []error
		for _, err := range res.Errors {
			errs = append(errs, errors.New(err))
		}

		context := contextHandle.RecoverContextAndDeferCleanup(cs.ItemID(dslOrigin.ItemID))
		return handler(nodeIds, res.Valid, errs, res.AllOk, context)
	})
}
