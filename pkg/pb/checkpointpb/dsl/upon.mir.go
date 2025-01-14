package checkpointpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/checkpointpb/types"
	commonpb "github.com/filecoin-project/mir/pkg/pb/commonpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types2 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Checkpoint](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponStableCheckpoint(m dsl.Module, handler func(sn types2.SeqNr, snapshot *commonpb.StateSnapshot, cert map[string][]uint8) error) {
	UponEvent[*types.Event_StableCheckpoint](m, func(ev *types.StableCheckpoint) error {
		return handler(ev.Sn, ev.Snapshot, ev.Cert)
	})
}
