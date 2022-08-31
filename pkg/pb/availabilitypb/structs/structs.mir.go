package availabilitypbstructs

import (
	availabilitypb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	structs "github.com/filecoin-project/mir/pkg/pb/contextstorepb/structs"
	dslpb "github.com/filecoin-project/mir/pkg/pb/dslpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	isEvent_Type()
	Pb() availabilitypb.Event_Type
}

func Event_TypeFromPb(pb availabilitypb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *availabilitypb.Event_RequestCert:
		return &Event_RequestCert{RequestCert: pb.RequestCert}
	case *availabilitypb.Event_NewCert:
		return &Event_NewCert{NewCert: pb.NewCert}
	case *availabilitypb.Event_VerifyCert:
		return &Event_VerifyCert{VerifyCert: pb.VerifyCert}
	case *availabilitypb.Event_CertVerified:
		return &Event_CertVerified{CertVerified: CertVerifiedFromPb(pb.CertVerified)}
	case *availabilitypb.Event_RequestTransactions:
		return &Event_RequestTransactions{RequestTransactions: pb.RequestTransactions}
	case *availabilitypb.Event_ProvideTransactions:
		return &Event_ProvideTransactions{ProvideTransactions: pb.ProvideTransactions}
	}
	return nil
}

type Event_RequestCert struct {
	RequestCert *availabilitypb.RequestCert
}

func (*Event_RequestCert) isEvent_Type() {}

func (w *Event_RequestCert) Pb() availabilitypb.Event_Type {
	return &availabilitypb.Event_RequestCert{RequestCert: w.RequestCert}
}

type Event_NewCert struct {
	NewCert *availabilitypb.NewCert
}

func (*Event_NewCert) isEvent_Type() {}

func (w *Event_NewCert) Pb() availabilitypb.Event_Type {
	return &availabilitypb.Event_NewCert{NewCert: w.NewCert}
}

type Event_VerifyCert struct {
	VerifyCert *availabilitypb.VerifyCert
}

func (*Event_VerifyCert) isEvent_Type() {}

func (w *Event_VerifyCert) Pb() availabilitypb.Event_Type {
	return &availabilitypb.Event_VerifyCert{VerifyCert: w.VerifyCert}
}

type Event_CertVerified struct {
	CertVerified *CertVerified
}

func (*Event_CertVerified) isEvent_Type() {}

func (w *Event_CertVerified) Pb() availabilitypb.Event_Type {
	return &availabilitypb.Event_CertVerified{CertVerified: w.CertVerified.Pb()}
}

type Event_RequestTransactions struct {
	RequestTransactions *availabilitypb.RequestTransactions
}

func (*Event_RequestTransactions) isEvent_Type() {}

func (w *Event_RequestTransactions) Pb() availabilitypb.Event_Type {
	return &availabilitypb.Event_RequestTransactions{RequestTransactions: w.RequestTransactions}
}

type Event_ProvideTransactions struct {
	ProvideTransactions *availabilitypb.ProvideTransactions
}

func (*Event_ProvideTransactions) isEvent_Type() {}

func (w *Event_ProvideTransactions) Pb() availabilitypb.Event_Type {
	return &availabilitypb.Event_ProvideTransactions{ProvideTransactions: w.ProvideTransactions}
}

func NewEvent(type_ Event_Type) *Event {
	return &Event{
		Type: type_,
	}
}

func (m *Event) Pb() *availabilitypb.Event {
	return &availabilitypb.Event{
		Type: m.Type.Pb(),
	}
}

func EventFromPb(pb *availabilitypb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

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

func RequestCertOrigin_TypeFromPb(pb availabilitypb.RequestCertOrigin_Type) RequestCertOrigin_Type {
	switch pb := pb.(type) {
	case *availabilitypb.RequestCertOrigin_ContextStore:
		return &RequestCertOrigin_ContextStore{ContextStore: structs.OriginFromPb(pb.ContextStore)}
	case *availabilitypb.RequestCertOrigin_Dsl:
		return &RequestCertOrigin_Dsl{Dsl: pb.Dsl}
	}
	return nil
}

type RequestCertOrigin_ContextStore struct {
	ContextStore *structs.Origin
}

func (*RequestCertOrigin_ContextStore) isRequestCertOrigin_Type() {}

func (w *RequestCertOrigin_ContextStore) Pb() availabilitypb.RequestCertOrigin_Type {
	return &availabilitypb.RequestCertOrigin_ContextStore{ContextStore: w.ContextStore.Pb()}
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

func VerifyCertOrigin_TypeFromPb(pb availabilitypb.VerifyCertOrigin_Type) VerifyCertOrigin_Type {
	switch pb := pb.(type) {
	case *availabilitypb.VerifyCertOrigin_ContextStore:
		return &VerifyCertOrigin_ContextStore{ContextStore: structs.OriginFromPb(pb.ContextStore)}
	case *availabilitypb.VerifyCertOrigin_Dsl:
		return &VerifyCertOrigin_Dsl{Dsl: pb.Dsl}
	}
	return nil
}

type VerifyCertOrigin_ContextStore struct {
	ContextStore *structs.Origin
}

func (*VerifyCertOrigin_ContextStore) isVerifyCertOrigin_Type() {}

func (w *VerifyCertOrigin_ContextStore) Pb() availabilitypb.VerifyCertOrigin_Type {
	return &availabilitypb.VerifyCertOrigin_ContextStore{ContextStore: w.ContextStore.Pb()}
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
