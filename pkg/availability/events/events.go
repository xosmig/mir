package events

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

func RequestCert(dest t.ModuleID, origin *apb.RequestCertOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_RequestCert{
			RequestCert: &apb.RequestCert{
				Origin: origin,
			},
		},
	})
}

func NewCert(dest t.ModuleID, cert *apb.Cert, origin *apb.RequestCertOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_NewCert{
			NewCert: &apb.NewCert{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

func RequestTransactions(dest t.ModuleID, cert *apb.Cert, origin *apb.RequestTransactionsOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_RequestTransactions{
			RequestTransactions: &apb.RequestTransactions{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

func ProvideTransactions(dest t.ModuleID, txs [][]byte, origin *apb.RequestTransactionsOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_ProvideTransactions{
			ProvideTransactions: &apb.ProvideTransactions{
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}

func VerifyCert(dest t.ModuleID, cert *apb.Cert, origin *apb.VerifyCertOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_VerifyCert{
			VerifyCert: &apb.VerifyCert{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

func CertVerified(dest t.ModuleID, valid bool, origin *apb.VerifyCertOrigin) *eventpb.Event {
	return Event(dest, &apb.Event{
		Type: &apb.Event_CertVerified{
			CertVerified: &apb.CertVerified{
				Valid:  valid,
				Origin: origin,
			},
		},
	})
}
