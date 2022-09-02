package availabilitypbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/availabilitypb/events"
	types "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
)

// Module-specific dsl functions for emitting events.

func CertVerified(m dsl.Module, destModule string, valid bool, err string, origin *types.VerifyCertOrigin) {
	dsl.EmitEvent(m, events.CertVerified(destModule, valid, err, origin))
}
