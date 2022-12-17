package dslpingpongpbevents

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func PingTime(destModule types.ModuleID) *types1.Event {
	return &types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Dslpingpong{
			Dslpingpong: &types2.Event{
				Type: &types2.Event_PingTime{
					PingTime: &types2.PingTime{},
				},
			},
		},
	}
}
