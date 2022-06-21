package dsl

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/pkg/errors"
)

// EventType represents a set of types assignable to the Type field of eventpb.Event.
// Copied from pkg/pb/eventpb/eventpb.pb.go.
// See https://github.com/golang/protobuf/issues/261 to know there is no nicer way to do this.
// TODO: consider replacing with a protoc plugin that would export such an interface.
type EventType interface {
	eventpb.Event_Init |
		eventpb.Event_Tick |
		eventpb.Event_WalAppend |
		eventpb.Event_WalEntry |
		eventpb.Event_WalTruncate |
		eventpb.Event_WalLoadAll |
		eventpb.Event_Request |
		eventpb.Event_HashRequest |
		eventpb.Event_HashResult |
		eventpb.Event_SignRequest |
		eventpb.Event_SignResult |
		eventpb.Event_VerifyNodeSigs |
		eventpb.Event_NodeSigsVerified |
		eventpb.Event_RequestReady |
		eventpb.Event_SendMessage |
		eventpb.Event_MessageReceived |
		eventpb.Event_Deliver |
		eventpb.Event_Iss |
		eventpb.Event_VerifyRequestSig |
		eventpb.Event_RequestSigVerified |
		eventpb.Event_StoreVerifiedRequest |
		eventpb.Event_AppSnapshotRequest |
		eventpb.Event_AppSnapshot |
		eventpb.Event_AppRestoreState |
		eventpb.Event_TimerDelay |
		eventpb.Event_TimerRepeat |
		eventpb.Event_TimerGarbageCollect |
		eventpb.Event_Bcb
}

// Dsl functions for emitting events

func SendMessage(m Module, destModule t.ModuleID, msg *messagepb.Message, dest []t.NodeID) {
	EmitEvent(m, events.SendMessage(destModule, msg, dest))
}

func SignRequest[C any](m Module, destModule t.ModuleID, data [][]byte, context C) {
	contextID := m.GetDslHandle().StoreContext(context)

	origin := &eventpb.SignOrigin{
		Module: m.GetModuleID().Pb(),
		Type: &eventpb.SignOrigin_Dsl{
			Dsl: &eventpb.DslOrigin{
				ContextID: contextID.Pb(),
			},
		},
	}
	EmitEvent(m, events.SignRequest(destModule, data, origin))
}

func VerifyOneNodeSig[C any](
	m Module,
	destModule t.ModuleID,
	data [][]byte,
	signature []byte,
	nodeID t.NodeID,
	context C,
) {
	VerifyNodeSigs(m, destModule, [][][]byte{data}, [][]byte{signature}, []t.NodeID{nodeID}, context)
}

// VerifyNodeSigs emits a signature verification event for a batch of signatures.
func VerifyNodeSigs[C any](
	m Module,
	destModule t.ModuleID,
	data [][][]byte,
	signatures [][]byte,
	nodeIDs []t.NodeID,
	context C,
) {
	contextID := m.GetDslHandle().StoreContext(context)

	origin := &eventpb.SigVerOrigin{
		Module: m.GetModuleID().Pb(),
		Type: &eventpb.SigVerOrigin_Dsl{
			Dsl: &eventpb.DslOrigin{
				ContextID: contextID.Pb(),
			},
		},
	}

	EmitEvent(m, events.VerifyNodeSigs(destModule, data, signatures, nodeIDs, origin))
}

func HashRequest[C any](m Module, destModule t.ModuleID, data [][][]byte, context C) {
	contextID := m.GetDslHandle().StoreContext(context)

	origin := &eventpb.HashOrigin{
		Module: m.GetModuleID().Pb(),
		Type: &eventpb.HashOrigin_Dsl{
			Dsl: &eventpb.DslOrigin{
				ContextID: contextID.Pb(),
			},
		},
	}

	EmitEvent(m, events.HashRequest(destModule, data, origin))
}

// Dsl functions for processing events

func UponRequest(m Module, handler func(clientId t.ClientID, reqNo uint64, data []byte, authenticator []byte) error) {
	RegisterEventHandler(m, func(eventType *eventpb.Event_Request) error {
		req := eventType.Request
		return handler(t.ClientID(req.ClientId), req.ReqNo, req.Data, req.Authenticator)
	})
}

func UponSignResult[C any](m Module, handler func(signature []byte, context C) error) {
	RegisterEventHandler(m, func(eventType *eventpb.Event_SignResult) error {
		res := eventType.SignResult

		dslOriginWrapper, ok := res.Origin.Type.(*eventpb.SignOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.GetDslHandle().RecoverAndCleanupContext(ContextID(dslOriginWrapper.Dsl.ContextID))
		context, ok := contextRaw.(C)
		if !ok {
			return nil
		}

		return handler(res.Signature, context)
	})
}

func UponNodeSigsVerified[C any](
	m Module,
	handler func(nodeIDs []t.NodeID, valid []bool, errs []error, allOK bool, context C) error,
) {
	RegisterEventHandler(m, func(eventType *eventpb.Event_NodeSigsVerified) error {
		res := eventType.NodeSigsVerified

		dslOriginWrapper, ok := res.Origin.Type.(*eventpb.SigVerOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.GetDslHandle().RecoverAndCleanupContext(ContextID(dslOriginWrapper.Dsl.ContextID))
		context, ok := contextRaw.(C)
		if !ok {
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

		return handler(nodeIds, res.Valid, errs, res.AllOk, context)
	})
}

func UponOneNodeSigVerified[C any](m Module, handler func(nodeID t.NodeID, valid bool, err error, context C) error) {
	UponNodeSigsVerified(m, func(nodeIDs []t.NodeID, valid []bool, errs []error, allOK bool, context C) error {
		for i := range nodeIDs {
			err := handler(nodeIDs[i], valid[i], errs[i], context)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func UponMessageReceived(m Module, handler func(from t.NodeID, msg *messagepb.Message) error) {
	RegisterEventHandler(m, func(eventType *eventpb.Event_MessageReceived) error {
		ev := eventType.MessageReceived
		return handler(t.NodeID(ev.From), ev.Msg)
	})
}
