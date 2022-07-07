package events

import (
	contextstore "github.com/filecoin-project/mir/pkg/contextstore"
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
)

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
