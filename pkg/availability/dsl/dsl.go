package dsl

import (
	aevents "github.com/filecoin-project/mir/pkg/availability/events"
	"github.com/filecoin-project/mir/pkg/dsl"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func RequestCert[C any](m dsl.Module, dest t.ModuleID, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &apb.RequestCertOrigin{
		Module: m.ModuleID().Pb(),
		Type: &apb.RequestCertOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, aevents.RequestCert(dest, origin))
}

func NewCert(m dsl.Module, dest t.ModuleID, cert *apb.Cert, origin *apb.RequestCertOrigin) {
	dsl.EmitEvent(m, aevents.NewCert(dest, cert, origin))
}

func RequestTransactions[C any](m dsl.Module, dest t.ModuleID, cert *apb.Cert, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &apb.RequestTransactionsOrigin{
		Module: dest.Pb(),
		Type: &apb.RequestTransactionsOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, aevents.RequestTransactions(dest, cert, origin))
}

func ProvideTransactions(m dsl.Module, dest t.ModuleID, txs [][]byte, origin *apb.RequestTransactionsOrigin) {
	dsl.EmitEvent(m, aevents.ProvideTransactions(dest, txs, origin))
}

func VerifyCert[C any](m dsl.Module, dest t.ModuleID, cert *apb.Cert, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &apb.VerifyCertOrigin{
		Module: dest.Pb(),
		Type: &apb.VerifyCertOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, aevents.VerifyCert(dest, cert, origin))
}

func CertVerified(m dsl.Module, dest t.ModuleID, err error, origin *apb.VerifyCertOrigin) {
	dsl.EmitEvent(m, aevents.CertVerified(dest, err, origin))
}

// Module-specific dsl functions for processing events.

func UponEvent[EvWrapper apb.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponEvent[*eventpb.Event_Availability](m, func(ev *apb.Event) error {
		evWrapper, ok := ev.Type.(EvWrapper)
		if !ok {
			return nil
		}
		return handler(evWrapper.Unwrap())
	})
}

func UponRequestCert(m dsl.Module, handler func(origin *apb.RequestCertOrigin) error) {
	UponEvent[*apb.Event_RequestCert](m, func(ev *apb.RequestCert) error {
		return handler(ev.Origin)
	})
}

func UponNewCert[C any](m dsl.Module, handler func(context *C) error) {
	UponEvent[*apb.Event_NewCert](m, func(ev *apb.NewCert) error {
		OriginWrapper, ok := ev.Origin.Type.(*apb.RequestCertOrigin_Dsl)
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

func UponRequestTransactions(m dsl.Module, handler func(cert *apb.Cert, origin *apb.RequestTransactionsOrigin) error) {
	UponEvent[*apb.Event_RequestTransactions](m, func(ev *apb.RequestTransactions) error {
		return handler(ev.Cert, ev.Origin)
	})
}

func UponProvideTransactions[C any](m dsl.Module, handler func(txs [][]byte, context *C) error) {
	UponEvent[*apb.Event_ProvideTransactions](m, func(ev *apb.ProvideTransactions) error {
		originWrapper, ok := ev.Origin.Type.(*apb.RequestTransactionsOrigin_Dsl)
		if !ok {
			return nil
		}

		contextRaw := m.DslHandle().RecoverAndCleanupContext(dsl.ContextID(originWrapper.Dsl.ContextID))
		context, ok := contextRaw.(*C)
		if !ok {
			return nil
		}

		return handler(ev.Txs, context)
	})
}

func UponVerifyCert(m dsl.Module, handler func(cert *apb.Cert, origin *apb.VerifyCertOrigin) error) {
	UponEvent[*apb.Event_VerifyCert](m, func(ev *apb.VerifyCert) error {
		return handler(ev.Cert, ev.Origin)
	})
}

func UponCertVerified[C any](m dsl.Module, handler func(err error, context *C) error) {
	UponEvent[*apb.Event_CertVerified](m, func(ev *apb.CertVerified) error {
		originWrapper, ok := ev.Origin.Type.(*apb.VerifyCertOrigin_Dsl)
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
