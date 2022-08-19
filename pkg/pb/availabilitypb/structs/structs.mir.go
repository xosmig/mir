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

func (m *CertVerified) ToPb() *availabilitypb.CertVerified {
	return &availabilitypb.CertVerified{Valid: m.Valid, Err: m.Err, Origin: m.Origin.ToPb()}
}

type RequestCertOrigin struct {
	Module types.ModuleID
}

func (m *RequestCertOrigin) ToPb() *availabilitypb.RequestCertOrigin {
	return &availabilitypb.RequestCertOrigin{Module: (string)(m.Module)}
}

type VerifyCertOrigin struct {
	Module string
}

func (m *VerifyCertOrigin) ToPb() *availabilitypb.VerifyCertOrigin {
	return &availabilitypb.VerifyCertOrigin{Module: m.Module}
}
