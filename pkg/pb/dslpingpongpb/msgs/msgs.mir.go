package dslpingpongpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func Ping(destModule types.ModuleID, seqNr types.SeqNr) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Dslpingpong{
			Dslpingpong: &types2.Message{
				Type: &types2.Message_Ping{
					Ping: &types2.Ping{
						SeqNr: seqNr,
					},
				},
			},
		},
	}
}

func Pong(destModule types.ModuleID, seqNr types.SeqNr) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Dslpingpong{
			Dslpingpong: &types2.Message{
				Type: &types2.Message_Pong{
					Pong: &types2.Pong{
						SeqNr: seqNr,
					},
				},
			},
		},
	}
}
