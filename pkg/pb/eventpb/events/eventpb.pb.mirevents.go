package eventpbevents

import (
	contextstore "github.com/filecoin-project/mir/pkg/contextstore"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

// convertSlice is an auxiliary functions to wrap / unwrap slices.
func convertSlice[T, R any](ts []T, f func(T) R) []R {
	rs := make([]R, len(ts))
	for i := range ts {
		rs[i] = f(ts[i])

	}
	return rs
}

func Init() *eventpb.Init {
	return &eventpb.Init{}
}

func ContextStoreOrigin(itemID contextstore.ItemID) *eventpb.ContextStoreOrigin {
	return &eventpb.ContextStoreOrigin{
		ItemID: uint64(itemID),
	}
}

func DslOrigin(contextID dsl.ContextID) *eventpb.DslOrigin {
	return &eventpb.DslOrigin{
		ContextID: uint64(contextID),
	}
}

func SignResult(signature []byte, origin *eventpb.SignOrigin) *eventpb.SignResult {
	return &eventpb.SignResult{
		Signature: signature,
		Origin:    origin,
	}
}

func SigVerData(data [][]byte) *eventpb.SigVerData {
	return &eventpb.SigVerData{
		Data: data,
	}
}

func VerifyNodeSigs(data []*eventpb.SigVerData, signatures [][]byte, origin *eventpb.SigVerOrigin, nodeIds []types.NodeID) *eventpb.VerifyNodeSigs {
	return &eventpb.VerifyNodeSigs{
		Data:       data,
		Signatures: signatures,
		Origin:     origin,
		NodeIds:    convertSlice(nodeIds, func(t types.NodeID) string { return string(t) }),
	}
}
