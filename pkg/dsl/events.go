package dsl

import (
	"fmt"
	cs "github.com/filecoin-project/mir/pkg/contextstore"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/pkg/errors"
)

//

// Event-specific dsl functions for emitting events (wrappers around the functions defined in pkg/events)

func SendMessage(m Module, destModule t.ModuleID, msg *messagepb.Message, dest []t.NodeID) {
	EmitEvent(m, events.SendMessage(destModule, msg, dest))
}

func SignRequest[T any](m Module, destModule t.ModuleID, data [][]byte, contextStore cs.ContextStore[T], context T) {
	//impl := m.GetDslHandle().impl

	if !contextStore.HasID() {
		panic(fmt.Errorf("no handlers are assigned to this contextstore of type '%T'", contextStore))
	}

	itemID := contextStore.Store(context)

	origin := &eventpb.SignOrigin{
		Module: m.GetDslHandle().impl.moduleID.Pb(),
		Type: &eventpb.SignOrigin_Dsl{
			Dsl: dslOrigin(contextStore.GetID(), itemID),
		},
	}
	EmitEvent(m, events.SignRequest(destModule, data, origin))
}

// VerifyNodeSigs emits a signature verification event for a batch of signatures.
func VerifyNodeSigs(
	m Module,
	destModule t.ModuleID,
	data [][][]byte,
	signatures [][]byte,
	nodeIDs []t.NodeID,
	origin *eventpb.SigVerOrigin,
) {
	EmitEvent(m, events.VerifyNodeSigs(destModule, data, signatures, nodeIDs, origin))
}

// Event-specific dsl functions for processing events

func UponRequest(m Module, handler func(clientId t.ClientID, reqNo uint64, data []byte, authenticator []byte) error) {
	UponEvent[eventpb.Event_Request](m, func(req *requestpb.Request) error {
		return handler(t.ClientID(req.ClientId), req.ReqNo, req.Data, req.Authenticator)
	})
}

func UponSignResult[T any](m Module, contextStore cs.ContextStore[T], handler func(signature []byte, context T) error) {
	storeID := m.GetDslHandle().impl.getContextStoreId(contextStore)

	UponEvent[eventpb.Event_SignRequest](m, func(res *eventpb.SignResult) error {
		csOrigin, ok := res.Origin.Type.(*eventpb.SignOrigin_ContextStore)
		if !ok {
			return nil
		}

		if csOrigin.ContextStore.StoreID != storeID {
			return nil
		}

		return handler(res.Signature, contextStore.RecoverAndDispose(cs.ItemID(csOrigin.ContextStore.ItemID)))
	})
}

func UponNodeSigsVerified[T any](
	m Module,
	contextStore cs.ContextStore[T],
	handler func(nodeIds []t.NodeID, valid []bool, errs []error, allOk bool, context T) error) {
	UponEvent[eventpb.Event_NodeSigsVerified](m, func(res *eventpb.NodeSigsVerified) error {
		csOrigin, ok := res.Origin.Type.(*eventpb.SigVerOrigin_ContextStore)
		if !ok {
			return nil
		}

		if cs.StoreID(csOrigin.ContextStore.StoreID) != contextStore.Id() {
			return nil
		}

		if cs.StoreID(csOrigin.ContextStore.StoreID) != contextStore.Id() {
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

		context := contextStore.RecoverAndDispose(cs.ItemID(csOrigin.ContextStore.ItemID))
		return handler(nodeIds, res.Valid, errs, res.AllOk, context)
	})
}

func UponMessageReceived(m Module, handler func(from t.NodeID, msg *messagepb.Message) error) {
	UponEvent[eventpb.Event_MessageReceived](m, func(ev *eventpb.MessageReceived) error {
		return handler(t.NodeID(ev.From), ev.Msg)
	})
}
