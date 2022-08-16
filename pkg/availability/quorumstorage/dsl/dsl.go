package dsl

import (
	qsevents "github.com/filecoin-project/mir/pkg/availability/quorumstorage/events"
	"github.com/filecoin-project/mir/pkg/dsl"
	qspb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/quorumstoragepb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func StoreBatchOnQuorum[C any](m dsl.Module, dest t.ModuleID, metadata []byte, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &qspb.StoreBatchOnQuorumOrigin{
		Module: m.ModuleID().Pb(),
		Type: &qspb.StoreBatchOnQuorumOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, qsevents.StoreBatchOnQuorum(dest, metadata, origin))
}

func BatchStoredOnQuorum(m dsl.Module, dest t.ModuleID, cert *qspb.RetrieveCert, origin *qspb.StoreBatchOnQuorumOrigin) {
	dsl.EmitEvent(m, qsevents.BatchStoredOnQuorum(dest, cert, origin))
}

func VerifyBatch[C any](m dsl.Module, dest t.ModuleID, metadata []byte, origin *qspb.VerifyBatchOrigin) {
	dsl.EmitEvent(m, qsevents.VerifyBatch(dest, metadata, origin))
}

func BatchVerified(m dsl.Module, dest t.ModuleID, err error, origin *qspb.VerifyBatchOrigin) {
	dsl.EmitEvent(m, qsevents.BatchVerified(dest, err, origin))
}

func VerifyCert[C any](m dsl.Module, dest t.ModuleID, cert *qspb.RetrieveCert, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &qspb.VerifyCertOrigin{
		Module: dest.Pb(),
		Type: &qspb.VerifyCertOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, qsevents.VerifyCert(dest, cert, origin))
}

func CertVerified(m dsl.Module, dest t.ModuleID, err error, origin *qspb.VerifyCertOrigin) {
	dsl.EmitEvent(m, qsevents.CertVerified(dest, err, origin))
}

func RetrieveTransactions[C any](m dsl.Module, dest t.ModuleID, cert *qspb.RetrieveCert, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &qspb.RetrieveTransactionsOrigin{
		Module: dest.Pb(),
		Type: &qspb.RetrieveTransactionsOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, qsevents.RetrieveTransactions(dest, cert, origin))
}

func TransactionsRetrieved(m dsl.Module, dest t.ModuleID, txIDs []t.TxID, txs []*requestpb.Request, origin *qspb.RetrieveTransactionsOrigin) {
	dsl.EmitEvent(m, qsevents.TransactionsRetrieved(dest, txIDs, txs, origin))
}

// Module-specific dsl functions for processing events.

// UponEvent registers a handler for the given batch storage event type.
func UponEvent[EvWrapper qspb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_QuorumStorage](m, func(ev *qspb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}

// UponStoreBatchOnQuorum registers a handler for the StoreBatchOnQuorum event.
func UponStoreBatchOnQuorum(m dsl.Module, handler func(origin *qspb.StoreBatchOnQuorumOrigin) error) {
	UponEvent[*qspb.Event_Store](m, func(ev *qspb.StoreBatchOnQuorum) error {
		return handler(ev.Origin)
	})
}

// UponBatchStoredOnQuorum registers a handler for the BatchStoredOnQuorum event.
func UponBatchStoredOnQuorum[C any](m dsl.Module, handler func(context *C) error) {
	UponEvent[*qspb.Event_Stored](m, func(ev *qspb.BatchStoredOnQuorum) error {
		OriginWrapper, ok := ev.Origin.Type.(*qspb.StoreBatchOnQuorumOrigin_Dsl)
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
func UponVerifyCert(m dsl.Module, handler func(cert *qspb.RetrieveCert, origin *qspb.VerifyCertOrigin) error) {
	UponEvent[*qspb.Event_VerifyCert](m, func(ev *qspb.VerifyCert) error {
		return handler(ev.Cert, ev.Origin)
	})
}

// UponCertVerified registers a handler for the CertVerified events.
func UponCertVerified[C any](m dsl.Module, handler func(err error, context *C) error) {
	UponEvent[*qspb.Event_CertVerified](m, func(ev *qspb.CertVerified) error {
		originWrapper, ok := ev.Origin.Type.(*qspb.VerifyCertOrigin_Dsl)
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
func UponRetrieveTransactions(m dsl.Module, handler func(cert *qspb.RetrieveCert, origin *qspb.RetrieveTransactionsOrigin) error) {
	UponEvent[*qspb.Event_RetrieveTransactions](m, func(ev *qspb.RetrieveTransactions) error {
		return handler(ev.Cert, ev.Origin)
	})
}

// UponTransactionsRetrieved registers a handler for the TransactionsRetrieved events.
func UponTransactionsRetrieved[C any](m dsl.Module, handler func(txIDs []t.TxID, txs []*requestpb.Request, context *C) error) {
	UponEvent[*qspb.Event_TransactionsRetrieved](m, func(ev *qspb.TransactionsRetrieved) error {
		originWrapper, ok := ev.Origin.Type.(*qspb.RetrieveTransactionsOrigin_Dsl)
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
