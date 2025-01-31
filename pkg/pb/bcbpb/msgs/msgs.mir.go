package bcbpbmsgs

import (
	types2 "github.com/filecoin-project/mir/pkg/pb/bcbpb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types "github.com/filecoin-project/mir/pkg/types"
)

func StartMessage(destModule types.ModuleID, data []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Bcb{
			Bcb: &types2.Message{
				Type: &types2.Message_StartMessage{
					StartMessage: &types2.StartMessage{
						Data: data,
					},
				},
			},
		},
	}
}

func EchoMessage(destModule types.ModuleID, signature []uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Bcb{
			Bcb: &types2.Message{
				Type: &types2.Message_EchoMessage{
					EchoMessage: &types2.EchoMessage{
						Signature: signature,
					},
				},
			},
		},
	}
}

func FinalMessage(destModule types.ModuleID, data []uint8, signers []types.NodeID, signatures [][]uint8) *types1.Message {
	return &types1.Message{
		DestModule: destModule,
		Type: &types1.Message_Bcb{
			Bcb: &types2.Message{
				Type: &types2.Message_FinalMessage{
					FinalMessage: &types2.FinalMessage{
						Data:       data,
						Signers:    signers,
						Signatures: signatures,
					},
				},
			},
		},
	}
}
