package bcbpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/bcbpb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Bcb](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponStartMessage(m dsl.Module, handler func(data []uint8) error) {
	UponMessageReceived[*types.Message_StartMessage](m, func(msg *types.StartMessage) error {
		return handler(msg.Data)
	})
}

func UponEchoMessage(m dsl.Module, handler func(signature []uint8) error) {
	UponMessageReceived[*types.Message_EchoMessage](m, func(msg *types.EchoMessage) error {
		return handler(msg.Signature)
	})
}

func UponFinalMessage(m dsl.Module, handler func(data []uint8, signers []string, signatures [][]uint8) error) {
	UponMessageReceived[*types.Message_FinalMessage](m, func(msg *types.FinalMessage) error {
		return handler(msg.Data, msg.Signers, msg.Signatures)
	})
}
