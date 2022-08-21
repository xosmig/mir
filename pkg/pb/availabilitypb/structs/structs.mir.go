package structs

import (
	availabilitypb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

type CertVerified struct {
	Valid  bool
	Err    string
	Origin *VerifyCertOrigin
}

func NewCertVerified(valid bool, err string, origin *VerifyCertOrigin) *CertVerified {
	return &CertVerified{
		Valid:  valid,
		Err:    err,
		Origin: origin,
	}
}

func (m *CertVerified) Pb() *availabilitypb.CertVerified {
	return &availabilitypb.CertVerified{
		Valid:  m.Valid,
		Err:    m.Err,
		Origin: m.Origin.Pb(),
	}
}

func CertVerifiedFromPb(pb *availabilitypb.CertVerified) *CertVerified {
	return &CertVerified{
		Valid:  pb.Valid,
		Err:    pb.Err,
		Origin: VerifyCertOriginFromPb(pb.Origin),
	}
}

type RequestCertOrigin struct {
	Module types.ModuleID
	Type   RequestCertOrigin_Type
}

type RequestCertOrigin_Type interface {
	isRequestCertOrigin_Type()
	Pb() availabilitypb.RequestCertOrigin_Type
}

type RequestCertOrigin_ContextStore struct {
	ContextStore *contextstorepb.Origin
}

func (*RequestCertOrigin_ContextStore) isRequestCertOrigin_Type() {}

func (w *RequestCertOrigin_ContextStore) Pb() availabilitypb.RequestCertOrigin_Type {
	return &availabilitypb.RequestCertOrigin_ContextStore{ContextStore: w.ContextStore}
}

type RequestCertOrigin_Dsl struct {
	Dsl *dslpb.Origin
}

func (*RequestCertOrigin_Dsl) isRequestCertOrigin_Type() {}

func (w *RequestCertOrigin_Dsl) Pb() availabilitypb.RequestCertOrigin_Type {
	return &availabilitypb.RequestCertOrigin_Dsl{Dsl: w.Dsl}
}

func NewRequestCertOrigin(module types.ModuleID, type_ RequestCertOrigin_Type) *RequestCertOrigin {
	return &RequestCertOrigin{
		Module: module,
		Type:   type_,
	}
}

func (m *RequestCertOrigin) Pb() *availabilitypb.RequestCertOrigin {
	return &availabilitypb.RequestCertOrigin{
		Module: (string)(m.Module),
		Type:   m.Type.Pb(),
	}
}

func RequestCertOriginFromPb(pb *availabilitypb.RequestCertOrigin) *RequestCertOrigin {
	return &RequestCertOrigin{
		Module: (types.ModuleID)(pb.Module),
		Type:   RequestCertOrigin_TypeFromPb(pb.Type),
	}
}

type VerifyCertOrigin struct {
	Module string
	Type   VerifyCertOrigin_Type
}

type VerifyCertOrigin_Type interface {
	isVerifyCertOrigin_Type()
	Pb() availabilitypb.VerifyCertOrigin_Type
}

type VerifyCertOrigin_ContextStore struct {
	ContextStore *contextstorepb.Origin
}

func (*VerifyCertOrigin_ContextStore) isVerifyCertOrigin_Type() {}

func (w *VerifyCertOrigin_ContextStore) Pb() availabilitypb.VerifyCertOrigin_Type {
	return &availabilitypb.VerifyCertOrigin_ContextStore{ContextStore: w.ContextStore}
}

type VerifyCertOrigin_Dsl struct {
	Dsl *dslpb.Origin
}

func (*VerifyCertOrigin_Dsl) isVerifyCertOrigin_Type() {}

func (w *VerifyCertOrigin_Dsl) Pb() availabilitypb.VerifyCertOrigin_Type {
	return &availabilitypb.VerifyCertOrigin_Dsl{Dsl: w.Dsl}
}

func NewVerifyCertOrigin(module string, type_ VerifyCertOrigin_Type) *VerifyCertOrigin {
	return &VerifyCertOrigin{
		Module: module,
		Type:   type_,
	}
}

func (m *VerifyCertOrigin) Pb() *availabilitypb.VerifyCertOrigin {
	return &availabilitypb.VerifyCertOrigin{
		Module: m.Module,
		Type:   m.Type.Pb(),
	}
}

func VerifyCertOriginFromPb(pb *availabilitypb.VerifyCertOrigin) *VerifyCertOrigin {
	return &VerifyCertOrigin{
		Module: pb.Module,
		Type:   VerifyCertOrigin_TypeFromPb(pb.Type),
	}
}
