package availabilitypb

import (
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
)

type RequestCertOrigin_Type = isRequestCertOrigin_Type

type RequestCertOrigin_TypeWrapper[Ev any] interface {
	RequestCertOrigin_Type
	Unwrap() *Ev
}

func (p *RequestCertOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return p.ContextStore
}

func (p *RequestCertOrigin_Dsl) Unwrap() *dslpb.Origin {
	return p.Dsl
}

type VerifyCertOrigin_Type = isVerifyCertOrigin_Type

type VerifyCertOrigin_TypeWrapper[Ev any] interface {
	VerifyCertOrigin_Type
	Unwrap() *Ev
}

func (p *VerifyCertOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return p.ContextStore
}

func (p *VerifyCertOrigin_Dsl) Unwrap() *dslpb.Origin {
	return p.Dsl
}
