package dsl

import (
	bsevents "github.com/filecoin-project/mir/pkg/batchstorage/events"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/pb/batchstoragepb"
	bspb "github.com/filecoin-project/mir/pkg/pb/batchstoragepb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func StoreBatchOnQuorum[C any](m dsl.Module, dest t.ModuleID, metadata []byte, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &batchstoragepb.StoreBatchOnQuorumOrigin{
		Module: m.ModuleID().Pb(),
		Type: &bspb.StoreBatchOnQuorumOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, bsevents.StoreBatchOnQuorum(dest, metadata, origin))
}

func BatchStoredOnQuorum(m dsl.Module, dest t.ModuleID, cert *bspb.RetrieveCert, origin *bspb.StoreBatchOnQuorumOrigin) {
	dsl.EmitEvent(m, bsevents.BatchStoredOnQuorum(dest, cert, origin))
}

func VerifyBatch[C any](m dsl.Module, dest t.ModuleID, metadata []byte, requestID uint64) {
	dsl.EmitEvent(m, bsevents.VerifyBatch(dest, metadata, requestID))
}

func BatchVerified(m dsl.Module, dest t.ModuleID, err error, requestID uint64) {
	dsl.EmitEvent(m, bsevents.BatchVerified(dest, err, requestID))
}

func VerifyCert[C any](m dsl.Module, dest t.ModuleID, cert *bspb.RetrieveCert, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &bspb.VerifyCertOrigin{
		Module: dest.Pb(),
		Type: &bspb.VerifyCertOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, bsevents.VerifyCert(dest, cert, origin))
}

func CertVerified(m dsl.Module, dest t.ModuleID, err error, origin *bspb.VerifyCertOrigin) {
	dsl.EmitEvent(m, bsevents.CertVerified(dest, err, origin))
}

func RetrieveTransactions[C any](m dsl.Module, dest t.ModuleID, cert *bspb.RetrieveCert, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &bspb.RetrieveTransactionsOrigin{
		Module: dest.Pb(),
		Type: &bspb.RetrieveTransactionsOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, bsevents.RetrieveTransactions(dest, cert, origin))
}

func TransactionsRetrieved(m dsl.Module, dest t.ModuleID, txIDs []t.TxID, txs []*requestpb.Request, origin *bspb.RetrieveTransactionsOrigin) {
	dsl.EmitEvent(m, bsevents.TransactionsRetrieved(dest, txIDs, txs, origin))
}

// Module-specific dsl functions for processing events.

// UponEvent registers a handler for the given batch storage event type.
func UponEvent[EvWrapper bspb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_BatchStorage](m, func(ev *bspb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}

// UponStoreBatchOnQuorum registers a handler for the StoreBatchOnQuorum event.
func UponStoreBatchOnQuorum(m dsl.Module, handler func(origin *bspb.StoreBatchOnQuorumOrigin) error) {
	UponEvent[*bspb.Event_Store](m, func(ev *bspb.StoreBatchOnQuorum) error {
		return handler(ev.Origin)
	})
}

// UponBatchStoredOnQuorum registers a handler for the BatchStoredOnQuorum event.
func UponBatchStoredOnQuorum[C any](m dsl.Module, handler func(context *C) error) {
	UponEvent[*bspb.Event_Stored](m, func(ev *bspb.BatchStoredOnQuorum) error {
		OriginWrapper, ok := ev.Origin.Type.(*bspb.StoreBatchOnQuorumOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(OriginWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(context)
	})
}

// UponVerifyCert registers a handler for the VerifyCert events.
func UponVerifyCert(m dsl.Module, handler func(cert *bspb.RetrieveCert, origin *bspb.VerifyCertOrigin) error) {
	UponEvent[*bspb.Event_VerifyCert](m, func(ev *bspb.VerifyCert) error {
		return handler(ev.Cert, ev.Origin)
	})
}

// UponCertVerified registers a handler for the CertVerified events.
func UponCertVerified[C any](m dsl.Module, handler func(err error, context *C) error) {
	UponEvent[*bspb.Event_CertVerified](m, func(ev *bspb.CertVerified) error {
		originWrapper, ok := ev.Origin.Type.(*bspb.VerifyCertOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(t.ErrorFromPb(ev.Valid, ev.Err), context)
	})
}

// UponRetrieveTransactions registers a handler for the RetrieveTransactions events.
func UponRetrieveTransactions(m dsl.Module, handler func(cert *bspb.RetrieveCert, origin *bspb.RetrieveTransactionsOrigin) error) {
	UponEvent[*bspb.Event_RetrieveTransactions](m, func(ev *bspb.RetrieveTransactions) error {
		return handler(ev.Cert, ev.Origin)
	})
}

// UponTransactionsRetrieved registers a handler for the TransactionsRetrieved events.
func UponTransactionsRetrieved[C any](m dsl.Module, handler func(txIDs []t.TxID, txs []*requestpb.Request, context *C) error) {
	UponEvent[*bspb.Event_TransactionsRetrieved](m, func(ev *bspb.TransactionsRetrieved) error {
		originWrapper, ok := ev.Origin.Type.(*bspb.RetrieveTransactionsOrigin_Dsl)
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
