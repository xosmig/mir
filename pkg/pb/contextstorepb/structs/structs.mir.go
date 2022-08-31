package contextstorepbstructs

import contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"

type Origin struct {
	ItemID uint64
}

func (m *Origin) Pb() *contextstorepb.Origin {
	return &contextstorepb.Origin{
		ItemID: m.ItemID,
	}
}

func OriginFromPb(pb *contextstorepb.Origin) *Origin {
	return &Origin{
		ItemID: pb.ItemID,
	}
}
