package availabilitypb

import (
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
	reflect "reflect"
)

type RequestCertOrigin_Type = isRequestCertOrigin_Type

func (*RequestCertOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*RequestCertOrigin_ContextStore)(nil)),
		reflect.TypeOf((*RequestCertOrigin_Dsl)(nil)),
	}
}

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

func (*VerifyCertOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*VerifyCertOrigin_ContextStore)(nil)),
		reflect.TypeOf((*VerifyCertOrigin_Dsl)(nil)),
	}
}

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
