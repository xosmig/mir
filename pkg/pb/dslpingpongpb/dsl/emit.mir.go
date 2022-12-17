package dslpingpongpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb/events"
	types "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for emitting events.

func PingTime(m dsl.Module, destModule types.ModuleID) {
	dsl.EmitMirEvent(m, events.PingTime(destModule))
}
