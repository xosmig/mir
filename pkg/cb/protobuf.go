package cb

import (
	"github.com/filecoin-project/mir/pkg/pb/cbpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Message(self t.ModuleID, msg *cbpb.CBMessage) *messagepb.Message {
	return &messagepb.Message{DestModule: self.Pb(), Type: &messagepb.Message_Cb{Cb: msg}}
}

func StartMessage(self t.ModuleID, data []byte) *messagepb.Message {
	return Message(self, &cbpb.CBMessage{
		Type: &cbpb.CBMessage_StartMessage{
			StartMessage: &cbpb.StartMessage{Data: data},
		},
	})
}

func EchoMessage(self t.ModuleID, data []byte) *messagepb.Message {
	return Message(self, &cbpb.CBMessage{
		Type: &cbpb.CBMessage_EchoMessage{
			EchoMessage: &cbpb.EchoMessage{Data: data},
		},
	})
}

func FinalMessage(self t.ModuleID, data []byte, signers []t.NodeID, signatures [][]byte) *messagepb.Message {
	return Message(self, &cbpb.CBMessage{
		Type: &cbpb.CBMessage_FinalMessage{
			FinalMessage: &cbpb.FinalMessage{
				Data: data,
				
				Signatures: signatures,
			},
		},
	})
}
