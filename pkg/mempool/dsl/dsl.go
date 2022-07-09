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

func NewBatch(m dsl.Module, dest t.ModuleID, txIDs []t.TxID, txs [][]byte, origin *mppb.RequestBatchOrigin) {
	dsl.EmitEvent(m, mpevents.NewBatch(dest, txIDs, txs, origin))
}

func RequestTransactionIDs[C any](m dsl.Module, dest t.ModuleID, txs [][]byte, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &mppb.RequestTransactionIDsOrigin{
		Type: &mppb.RequestTransactionIDsOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, mpevents.RequestTransactionIDs(dest, txs, origin))
}

func TransactionIDsResponse(m dsl.Module, dest t.ModuleID, txIDs []t.TxID, origin *mppb.RequestTransactionIDsOrigin) {
	dsl.EmitEvent(m, mpevents.TransactionIDsResponse(dest, txIDs, origin))
}

func RequestBatchID[C any](m dsl.Module, dest t.ModuleID, txIDs []t.TxID, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &mppb.RequestBatchIDOrigin{
		Type: &mppb.RequestBatchIDOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, mpevents.RequestBatchID(dest, txIDs, origin))
}

func BatchIDResponse(m dsl.Module, dest t.ModuleID, batchID t.BatchID, origin *mppb.RequestBatchIDOrigin) {
	dsl.EmitEvent(m, mpevents.BatchIDResponse(dest, batchID, origin))
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

func UponNewBatch[C any](m dsl.Module, handler func(txIDs []t.TxID, txs [][]byte, context *C) error) {
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

		return handler(t.TxIDSlice(ev.TxIds), ev.Txs, context)
	})
}

func UponRequestTransactionIDs[C any](m dsl.Module, handler func(txs [][]byte, origin *mppb.RequestTransactionIDsOrigin) error) {
	UponEvent[*mppb.Event_RequestTransactionIds](m, func(ev *mppb.RequestTransactionIDs) error {
		return handler(ev.Txs, ev.Origin)
	})
}

func UponTransactionIDsResponse[C any](m dsl.Module, handler func(txIDs []t.TxID, context *C) error) {
	UponEvent[*mppb.Event_TransactionIdsResponse](m, func(ev *mppb.TransactionIDsResponse) error {
		originWrapper, ok := ev.Origin.Type.(*mppb.RequestTransactionIDsOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(t.TxIDSlice(ev.TxIds), context)
	})
}

func UponRequestBatchID[C any](m dsl.Module, handler func(txIDs []t.TxID, origin *mppb.RequestBatchIDOrigin) error) {
	UponEvent[*mppb.Event_RequestBatchId](m, func(ev *mppb.RequestBatchID) error {
		return handler(t.TxIDSlice(ev.TxIds), ev.Origin)
	})
}

func UponBatchIDResponse[C any](m dsl.Module, handler func(batchID t.BatchID, context *C) error) {
	UponEvent[*mppb.Event_BatchIdResponse](m, func(ev *mppb.BatchIDResponse) error {
		originWrapper, ok := ev.Origin.Type.(*mppb.RequestBatchIDOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(t.BatchID(ev.BatchId), context)
	})
}
