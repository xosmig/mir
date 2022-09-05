package availabilitypb

import (
	contextstorepb "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
)

type RequestCertOrigin_Type = isRequestCertOrigin_Type

type RequestCertOrigin_TypeWrapper[T any] interface {
	RequestCertOrigin_Type
	Unwrap() *T
}

func (w *RequestCertOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *RequestCertOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

type VerifyCertOrigin_Type = isVerifyCertOrigin_Type

type VerifyCertOrigin_TypeWrapper[T any] interface {
	VerifyCertOrigin_Type
	Unwrap() *T
}

func (w *VerifyCertOrigin_ContextStore) Unwrap() *contextstorepb.Origin {
	return w.ContextStore
}

func (w *VerifyCertOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}
