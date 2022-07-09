package events

import (
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	mppb "github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Event(dest t.ModuleID, ev *mppb.Event) *eventpb.Event {
	return &eventpb.Event{
		DestModule: dest.Pb(),
		Type: &eventpb.Event_Mempool{
			Mempool: ev,
		},
	}
}

func RequestBatch(dest t.ModuleID, origin *mppb.RequestBatchOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_RequestBatch{
			RequestBatch: &mppb.RequestBatch{
				Origin: origin,
			},
		},
	})
}

func NewBatch(dest t.ModuleID, txIDs [][]byte, txs [][]byte, origin *mppb.RequestBatchOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_NewBatch{
			NewBatch: &mppb.NewBatch{
				TxIds:  txIDs,
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}
