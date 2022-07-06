package availabilitydsl

import (
	aevents "github.com/filecoin-project/mir/pkg/availability/availabilityevents"
	"github.com/filecoin-project/mir/pkg/dsl"
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func RequestBatch[C any](m dsl.Module, dest t.ModuleID, context *C) {
	contextID := m.DslHandle().StoreContext(context)

	origin := &apb.RequestBatchOrigin{
		Module: m.ModuleID().Pb(),
		Type: &apb.RequestBatchOrigin_Dsl{
			Dsl: dsl.Origin(contextID),
		},
	}

	dsl.EmitEvent(m, aevents.RequestBatch(dest, origin))
}

func NewBatch(m dsl.Module, dest t.ModuleID, batchID []byte, cert *apb.Cert, origin *apb.RequestBatchOrigin) {
	dsl.EmitEvent(m, aevents.NewBatch(dest, batchID, cert, origin))
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

func UponRequestBatch(m dsl.Module, handler func(origin *apb.RequestBatchOrigin) error) {
	UponEvent[*apb.Event_RequestBatch](m, func(ev *apb.RequestBatch) error {
		return handler(ev.Origin)
	})
}

func UponNewBatch[C any](m dsl.Module, handler func(context *C) error) {
	UponEvent[*apb.Event_NewBatch](m, func(ev *apb.NewBatch) error {
		OriginWrapper, ok := ev.Origin.Type.(*apb.RequestBatchOrigin_Dsl)
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
