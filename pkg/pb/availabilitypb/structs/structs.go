package structs

import (
	availabilitypb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	types "github.com/filecoin-project/mir/pkg/types"
)

type CertVerified struct {
	Valid  bool
	Err    string
	Origin *availabilitypb.VerifyCertOrigin
}

func CertVerifiedFromPb(certVerified *availabilitypb.CertVerified) *CertVerified {
	return &CertVerified{Valid: (bool)(certVerified.Valid), Err: (string)(certVerified.Err), Origin: (*availabilitypb.VerifyCertOrigin)(certVerified.Origin)}
}
func CertVerifiedToPb(certVerified CertVerified) *availabilitypb.CertVerified {
	return &availabilitypb.CertVerified{Valid: (bool)(certVerified.Valid), Err: (string)(certVerified.Err), Origin: (*availabilitypb.VerifyCertOrigin)(certVerified.Origin)}
}

type RequestCertOrigin struct {
	Module types.ModuleID
}

func RequestCertOriginFromPb(requestCertOrigin *availabilitypb.RequestCertOrigin) *RequestCertOrigin {
	return &RequestCertOrigin{Module: (types.ModuleID)(requestCertOrigin.Module)}
}
func RequestCertOriginToPb(requestCertOrigin RequestCertOrigin) *availabilitypb.RequestCertOrigin {
	return &availabilitypb.RequestCertOrigin{Module: (string)(requestCertOrigin.Module)}
}
