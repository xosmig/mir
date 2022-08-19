package structs

import (
	availabilitypb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	types "github.com/filecoin-project/mir/pkg/types"
)

type CertVerified struct {
	Valid  bool
	Err    string
	Origin *VerifyCertOrigin
}

func NewCertVerified(valid bool, err string, origin *VerifyCertOrigin) *CertVerified {
	return &CertVerified{Valid: valid, Err: err, Origin: origin}
}

func (m *CertVerified) Pb() *availabilitypb.CertVerified {
	return &availabilitypb.CertVerified{Valid: m.Valid, Err: m.Err, Origin: m.Origin.Pb()}
}

func CertVerifiedFromPb(pb *availabilitypb.CertVerified) *CertVerified {
	return &CertVerified{Valid: pb.Valid, Err: pb.Err, Origin: VerifyCertOriginFromPb(pb.Origin)}
}

type RequestCertOrigin struct {
	Module types.ModuleID
}

func NewRequestCertOrigin(module types.ModuleID) *RequestCertOrigin {
	return &RequestCertOrigin{Module: module}
}

func (m *RequestCertOrigin) Pb() *availabilitypb.RequestCertOrigin {
	return &availabilitypb.RequestCertOrigin{Module: (string)(m.Module)}
}

func RequestCertOriginFromPb(pb *availabilitypb.RequestCertOrigin) *RequestCertOrigin {
	return &RequestCertOrigin{Module: (types.ModuleID)(pb.Module)}
}

type VerifyCertOrigin struct {
	Module string
}

func NewVerifyCertOrigin(module string) *VerifyCertOrigin {
	return &VerifyCertOrigin{Module: module}
}

func (m *VerifyCertOrigin) Pb() *availabilitypb.VerifyCertOrigin {
	return &availabilitypb.VerifyCertOrigin{Module: m.Module}
}

func VerifyCertOriginFromPb(pb *availabilitypb.VerifyCertOrigin) *VerifyCertOrigin {
	return &VerifyCertOrigin{Module: pb.Module}
}
