package dsl

import (
	"github.com/filecoin-project/mir/pkg/dsl"
	mpevents "github.com/filecoin-project/mir/pkg/mempool/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	mppb "github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func RequestBatch[C any](m dsl.Module, dest t.ModuleID, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &mppb.RequestBatchOrigin{
		Type: &mppb.RequestBatchOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, mpevents.RequestBatch(dest, origin))
}

func NewBatch(m dsl.Module, dest t.ModuleID, txIDs [][]byte, txs [][]byte, origin *mppb.RequestBatchOrigin) {
	dsl.EmitEvent(m, mpevents.NewBatch(dest, txIDs, txs, origin))
}

// Module-specific dsl functions for processing events.

func UponEvent[EvWrapper mppb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_Mempool](m, func(ev *mppb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}

func UponRequestBatch[C any](m dsl.Module, handler func(origin *mppb.RequestBatchOrigin) error) {
	UponEvent[*mppb.Event_RequestBatch](m, func(ev *mppb.RequestBatch) error {
		return handler(ev.Origin)
	})
}

func UponNewBatch[C any](m dsl.Module, handler func(txIDs [][]byte, txs [][]byte, context *C) error) {
	UponEvent[*mppb.Event_NewBatch](m, func(ev *mppb.NewBatch) error {
		originWrapper, ok := ev.Origin.Type.(*mppb.RequestBatchOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(ev.TxIds, ev.Txs, context)
	})
}
