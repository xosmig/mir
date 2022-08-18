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
	return &CertVerified{
		Err:    (string)(certVerified.Err),
		Origin: (*availabilitypb.VerifyCertOrigin)(certVerified.Origin),
		Valid:  (bool)(certVerified.Valid),
	}
}

func CertVerifiedToPb(certVerified CertVerified) *availabilitypb.CertVerified {
	return &availabilitypb.CertVerified{
		Err:    (string)(certVerified.Err),
		Origin: (*availabilitypb.VerifyCertOrigin)(certVerified.Origin),
		Valid:  (bool)(certVerified.Valid),
	}
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
