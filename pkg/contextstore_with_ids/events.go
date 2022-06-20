package contextstore_with_ids

import (
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Origin returns a ContextStoreOrigin protobuf containing the given id.
func Origin(storeID StoreID, itemID ItemID) *eventpb.ContextStoreOrigin {
	return &eventpb.ContextStoreOrigin{
		StoreID: storeID.Pb(),
		ItemID:  itemID.Pb(),
	}
}

// SignOrigin returns a SignOrigin protobuf containing moduleID and contextstore_with_ids.Origin(itemID).
func SignOrigin(moduleID t.ModuleID, storeID StoreID, itemID ItemID) *eventpb.SignOrigin {
	return &eventpb.SignOrigin{
		Module: moduleID.Pb(),
		Type: &eventpb.SignOrigin_Contextstore{
			Contextstore: Origin(storeID, itemID),
		},
	}
}

// SigVerOrigin returns a SigVerOrigin protobuf containing moduleID and contextstore_with_ids.Origin(itemID).
func SigVerOrigin(moduleID t.ModuleID, storeID StoreID, itemID ItemID) *eventpb.SigVerOrigin {
	return &eventpb.SigVerOrigin{
		Module: moduleID.Pb(),
		Type: &eventpb.SigVerOrigin_Contextstore{
			Contextstore: Origin(storeID, itemID),
		},
	}
}

// HashOrigin returns a HashOrigin protobuf containing moduleID and contextstore_with_ids.Origin(itemID).
func HashOrigin(moduleID t.ModuleID, storeID StoreID, itemID ItemID) *eventpb.HashOrigin {
	return &eventpb.HashOrigin{
		Module: moduleID.Pb(),
		Type: &eventpb.HashOrigin_Contextstore{
			Contextstore: Origin(storeID, itemID),
		},
	}
}
