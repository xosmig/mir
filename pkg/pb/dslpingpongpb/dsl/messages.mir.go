package dslpingpongpbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb/types"
	dsl1 "github.com/filecoin-project/mir/pkg/pb/messagepb/dsl"
	types2 "github.com/filecoin-project/mir/pkg/pb/messagepb/types"
	types1 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing net messages.

func UponMessageReceived[W types.Message_TypeWrapper[M], M any](m dsl.Module, handler func(from types1.NodeID, msg *M) error) {
	dsl1.UponMessageReceived[*types2.Message_Dslpingpong](m, func(from types1.NodeID, msg *types.Message) error {
		w, ok := msg.Type.(W)
		if !ok {
			return nil
		}

		return handler(from, w.Unwrap())
	})
}

func UponPingReceived(m dsl.Module, handler func(from types1.NodeID, seqNr types1.SeqNr) error) {
	UponMessageReceived[*types.Message_Ping](m, func(from types1.NodeID, msg *types.Ping) error {
		return handler(from, msg.SeqNr)
	})
}

func UponPongReceived(m dsl.Module, handler func(from types1.NodeID, seqNr types1.SeqNr) error) {
	UponMessageReceived[*types.Message_Pong](m, func(from types1.NodeID, msg *types.Pong) error {
		return handler(from, msg.SeqNr)
	})
}
