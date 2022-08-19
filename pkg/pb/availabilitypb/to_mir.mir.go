package availabilitypb

import (
	structs "github.com/filecoin-project/mir/pkg/pb/availabilitypb/structs"
	types "github.com/filecoin-project/mir/pkg/types"
)

func (pb *CertVerified) ToMir() *structs.CertVerified {
	return &structs.CertVerified{Valid: pb.Valid, Err: pb.Err, Origin: pb.Origin.ToMir()}
}
func (pb *RequestCertOrigin) ToMir() *structs.RequestCertOrigin {
	return &structs.RequestCertOrigin{Module: (types.ModuleID)(pb.Module)}
}
func (pb *VerifyCertOrigin) ToMir() *structs.VerifyCertOrigin {
	return &structs.VerifyCertOrigin{Module: pb.Module}
}
