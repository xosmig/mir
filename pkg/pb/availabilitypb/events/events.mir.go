package availabilitypbevents

import (
	structs "github.com/filecoin-project/mir/pkg/pb/availabilitypb/structs"
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
	structs1 "github.com/filecoin-project/mir/pkg/pb/eventpb/structs"
)

func CertVerified(destModule string, valid bool, err string, origin *structs.VerifyCertOrigin) *eventpb.Event {
	return (&structs1.Event{
		DestModule: destModule,
		Type: &structs1.Event_Availability{
			Availability: &structs.Event{
				Type: &structs.Event_CertVerified{
					CertVerified: &structs.CertVerified{
						Valid:  valid,
						Err:    err,
						Origin: origin,
					},
				},
			},
		},
	}).Pb()
}
