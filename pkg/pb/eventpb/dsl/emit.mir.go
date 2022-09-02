package eventpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/eventpb/events"
	types "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func Event(m dsl.Module, type_ types.Event_Type, destModule types1.ModuleID) {
	dsl.EmitMirEvent(m, events.Event(type_, destModule))
}
