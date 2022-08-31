package availabilitypbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	events "github.com/filecoin-project/mir/pkg/pb/availabilitypb/events"
	structs "github.com/filecoin-project/mir/pkg/pb/availabilitypb/structs"
)

func CertVerified(m dsl.Module, destModule string, valid bool, err string, origin *structs.VerifyCertOrigin) {
	dsl.EmitEvent(m, events.CertVerified(destModule, valid, err, origin))
}
