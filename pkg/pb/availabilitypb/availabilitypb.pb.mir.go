package availabilitypb

import (
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
	reflect "reflect"
)

type Event_Type = isEvent_Type

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_RequestCert)(nil)),
		reflect.TypeOf((*Event_NewCert)(nil)),
		reflect.TypeOf((*Event_VerifyCert)(nil)),
		reflect.TypeOf((*Event_CertVerified)(nil)),
		reflect.TypeOf((*Event_RequestTransactions)(nil)),
		reflect.TypeOf((*Event_ProvideTransactions)(nil)),
	}
}

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_RequestCert) Unwrap() *RequestCert {
	return p.RequestCert
}

func (p *Event_NewCert) Unwrap() *NewCert {
	return p.NewCert
}

func (p *Event_VerifyCert) Unwrap() *VerifyCert {
	return p.VerifyCert
}

func (p *Event_CertVerified) Unwrap() *CertVerified {
	return p.CertVerified
}

func (p *Event_RequestTransactions) Unwrap() *RequestTransactions {
	return p.RequestTransactions
}

func (p *Event_ProvideTransactions) Unwrap() *ProvideTransactions {
	return p.ProvideTransactions
}

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
