package availabilitypbevents

import (
	types "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
)

func CertVerified(destModule string, valid bool, err string, origin *types.VerifyCertOrigin) *eventpb.Event {
	return (&types1.Event{
		DestModule: destModule,
		Type: &types1.Event_Availability{
			Availability: &types.Event{
				Type: &types.Event_CertVerified{
					CertVerified: &types.CertVerified{
						Valid:  valid,
						Err:    err,
						Origin: origin,
					},
				},
			},
		},
	}).Pb()
}
