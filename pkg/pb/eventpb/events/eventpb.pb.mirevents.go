package events

import (
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
)

func Init() *eventpb.Init {
	return &eventpb.Init{}
}

func ContextStoreOrigin(itemID uint64) *eventpb.ContextStoreOrigin {
	return &eventpb.ContextStoreOrigin{
		ItemID: itemID,
	}
}
