package bcbpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/proto-converter/mirreflect"
	bcbpb "github.com/filecoin-project/mir/pkg/pb/bcbpb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() bcbpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb bcbpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *bcbpb.Event_Request:
		return &Event_Request{Request: pb.Request}
	case *bcbpb.Event_Deliver:
		return &Event_Deliver{Deliver: pb.Deliver}
	}
	return nil
}

type Event_Request struct {
	Request *bcbpb.BroadcastRequest
}

func (*Event_Request) isEvent_Type() {}

func (w *Event_Request) Unwrap() *bcbpb.BroadcastRequest {
	return w.Request
}

func (w *Event_Request) Pb() bcbpb.Event_Type {
	return &bcbpb.Event_Request{Request: w.Request}
}

func (*Event_Request) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcbpb.Event_Request]()}
}

type Event_Deliver struct {
	Deliver *bcbpb.Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Unwrap() *bcbpb.Deliver {
	return w.Deliver
}

func (w *Event_Deliver) Pb() bcbpb.Event_Type {
	return &bcbpb.Event_Deliver{Deliver: w.Deliver}
}

func (*Event_Deliver) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcbpb.Event_Deliver]()}
}

func EventFromPb(pb *bcbpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *bcbpb.Event {
	return &bcbpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*bcbpb.Event]()}
}
