package availabilityevents

import (
	apb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Event(dest t.ModuleID, ev *apb.Event) *eventpb.Event {
	return &eventpb.Event{
		DestModule: dest.Pb(),
		Type: &eventpb.Event_Availability{
			Availability: ev,
		},
	}
}

func RequestBatch(dest t.ModuleID, origin *apb.RequestBatchOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_RequestBatch{
			RequestBatch: &apb.RequestBatch{
				Origin: origin,
			},
		},
	})
}

func NewBatch(dest t.ModuleID, batchID []byte, cert *apb.Cert, origin *apb.RequestBatchOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_NewBatch{
			NewBatch: &apb.NewBatch{
				BatchId: batchID,
				Cert:    cert,
				Origin:  origin,
			},
		},
	})
}
