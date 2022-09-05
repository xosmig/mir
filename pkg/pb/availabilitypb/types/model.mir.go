package availabilitypbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	availabilitypb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	types1 "github.com/filecoin-project/mir/pkg/pb/contextstorepb/types"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type CertVerified struct {
	Valid  bool
	Err    string
	Origin *VerifyCertOrigin
}

func CertVerifiedFromPb(pb *availabilitypb.CertVerified) *CertVerified {
	return &CertVerified{
		Valid:  pb.Valid,
		Err:    pb.Err,
		Origin: VerifyCertOriginFromPb(pb.Origin),
	}
}

func (m *CertVerified) Pb() *availabilitypb.CertVerified {
	return &availabilitypb.CertVerified{
		Valid:  m.Valid,
		Err:    m.Err,
		Origin: (m.Origin).Pb(),
	}
}

func (*CertVerified) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*availabilitypb.CertVerified]()}
}

type RequestCertOrigin struct {
	Module types.ModuleID
	Type   RequestCertOrigin_Type
}

type RequestCertOrigin_Type interface {
	mirreflect.GeneratedType
	isRequestCertOrigin_Type()
	Pb() availabilitypb.RequestCertOrigin_Type
}

type RequestCertOrigin_TypeWrapper[T any] interface {
	RequestCertOrigin_Type
	Unwrap() *T
}

func RequestCertOrigin_TypeFromPb(pb availabilitypb.RequestCertOrigin_Type) RequestCertOrigin_Type {
	switch pb := pb.(type) {
	case *availabilitypb.RequestCertOrigin_ContextStore:
		return &RequestCertOrigin_ContextStore{ContextStore: types1.OriginFromPb(pb.ContextStore)}
	case *availabilitypb.RequestCertOrigin_Dsl:
		return &RequestCertOrigin_Dsl{Dsl: pb.Dsl}
	}
	return nil
}

type RequestCertOrigin_ContextStore struct {
	ContextStore *types1.Origin
}

func (*RequestCertOrigin_ContextStore) isRequestCertOrigin_Type() {}

func (w *RequestCertOrigin_ContextStore) Unwrap() *types1.Origin {
	return w.ContextStore
}

func (w *RequestCertOrigin_ContextStore) Pb() availabilitypb.RequestCertOrigin_Type {
	return &availabilitypb.RequestCertOrigin_ContextStore{ContextStore: (w.ContextStore).Pb()}
}

func (*RequestCertOrigin_ContextStore) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*availabilitypb.RequestCertOrigin_ContextStore]()}
}

type RequestCertOrigin_Dsl struct {
	Dsl *dslpb.Origin
}

func (*RequestCertOrigin_Dsl) isRequestCertOrigin_Type() {}

func (w *RequestCertOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

func (w *RequestCertOrigin_Dsl) Pb() availabilitypb.RequestCertOrigin_Type {
	return &availabilitypb.RequestCertOrigin_Dsl{Dsl: w.Dsl}
}

func (*RequestCertOrigin_Dsl) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*availabilitypb.RequestCertOrigin_Dsl]()}
}

func RequestCertOriginFromPb(pb *availabilitypb.RequestCertOrigin) *RequestCertOrigin {
	return &RequestCertOrigin{
		Module: (types.ModuleID)(pb.Module),
		Type:   RequestCertOrigin_TypeFromPb(pb.Type),
	}
}

func (m *RequestCertOrigin) Pb() *availabilitypb.RequestCertOrigin {
	return &availabilitypb.RequestCertOrigin{
		Module: (string)(m.Module),
		Type:   (m.Type).Pb(),
	}
}

func (*RequestCertOrigin) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*availabilitypb.RequestCertOrigin]()}
}

type VerifyCertOrigin struct {
	Module string
	Type   VerifyCertOrigin_Type
}

type VerifyCertOrigin_Type interface {
	mirreflect.GeneratedType
	isVerifyCertOrigin_Type()
	Pb() availabilitypb.VerifyCertOrigin_Type
}

type VerifyCertOrigin_TypeWrapper[T any] interface {
	VerifyCertOrigin_Type
	Unwrap() *T
}

func VerifyCertOrigin_TypeFromPb(pb availabilitypb.VerifyCertOrigin_Type) VerifyCertOrigin_Type {
	switch pb := pb.(type) {
	case *availabilitypb.VerifyCertOrigin_ContextStore:
		return &VerifyCertOrigin_ContextStore{ContextStore: types1.OriginFromPb(pb.ContextStore)}
	case *availabilitypb.VerifyCertOrigin_Dsl:
		return &VerifyCertOrigin_Dsl{Dsl: pb.Dsl}
	}
	return nil
}

type VerifyCertOrigin_ContextStore struct {
	ContextStore *types1.Origin
}

func (*VerifyCertOrigin_ContextStore) isVerifyCertOrigin_Type() {}

func (w *VerifyCertOrigin_ContextStore) Unwrap() *types1.Origin {
	return w.ContextStore
}

func (w *VerifyCertOrigin_ContextStore) Pb() availabilitypb.VerifyCertOrigin_Type {
	return &availabilitypb.VerifyCertOrigin_ContextStore{ContextStore: (w.ContextStore).Pb()}
}

func (*VerifyCertOrigin_ContextStore) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*availabilitypb.VerifyCertOrigin_ContextStore]()}
}

type VerifyCertOrigin_Dsl struct {
	Dsl *dslpb.Origin
}

func (*VerifyCertOrigin_Dsl) isVerifyCertOrigin_Type() {}

func (w *VerifyCertOrigin_Dsl) Unwrap() *dslpb.Origin {
	return w.Dsl
}

func (w *VerifyCertOrigin_Dsl) Pb() availabilitypb.VerifyCertOrigin_Type {
	return &availabilitypb.VerifyCertOrigin_Dsl{Dsl: w.Dsl}
}

func (*VerifyCertOrigin_Dsl) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*availabilitypb.VerifyCertOrigin_Dsl]()}
}

func VerifyCertOriginFromPb(pb *availabilitypb.VerifyCertOrigin) *VerifyCertOrigin {
	return &VerifyCertOrigin{
		Module: pb.Module,
		Type:   VerifyCertOrigin_TypeFromPb(pb.Type),
	}
}

func (m *VerifyCertOrigin) Pb() *availabilitypb.VerifyCertOrigin {
	return &availabilitypb.VerifyCertOrigin{
		Module: m.Module,
		Type:   (m.Type).Pb(),
	}
}

func (*VerifyCertOrigin) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*availabilitypb.VerifyCertOrigin]()}
}
