package eventpbevents

import (
	types "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

func Event(type_ types.Event_Type, destModule types1.ModuleID) *types.Event {
	return &types.Event{
		Type:       type_,
		DestModule: destModule,
	}
}
