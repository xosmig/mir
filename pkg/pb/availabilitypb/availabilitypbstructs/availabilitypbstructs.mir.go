package availabilitypbstructs

import (
	availabilitypb "github.com/filecoin-project/mir/pkg/pb/availabilitypb"
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
	types "github.com/filecoin-project/mir/pkg/types"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	PbWrapper() availabilitypb.Event_Type
}

func NewEvent(type_ Event_Type) *Event {
	return &Event{
		Type: type_,
	}
}

func (m *Event) PbWrapper() *eventpb.Event_Availability {
	return &eventpb.Event_Availability{Availability: m.Pb()}
}

func (m *Event) Pb() *availabilitypb.Event {
	return &availabilitypb.Event{
		Type: m.Type.PbWrapper(),
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
	PbWrapper() availabilitypb.RequestCertOrigin_Type
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
		Type:   m.Type.PbWrapper(),
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
	PbWrapper() availabilitypb.VerifyCertOrigin_Type
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
		Type:   m.Type.PbWrapper(),
	}
}

func VerifyCertOriginFromPb(pb *availabilitypb.VerifyCertOrigin) *VerifyCertOrigin {
	return &VerifyCertOrigin{
		Module: pb.Module,
		Type:   VerifyCertOrigin_TypeFromPb(pb.Type),
	}
}
