package dslpingpongpbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	dslpingpongpb "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb"
	types "github.com/filecoin-project/mir/pkg/types"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Event struct {
	Type Event_Type
}

type Event_Type interface {
	mirreflect.GeneratedType
	isEvent_Type()
	Pb() dslpingpongpb.Event_Type
}

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func Event_TypeFromPb(pb dslpingpongpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *dslpingpongpb.Event_PingTime:
		return &Event_PingTime{PingTime: PingTimeFromPb(pb.PingTime)}
	}
	return nil
}

type Event_PingTime struct {
	PingTime *PingTime
}

func (*Event_PingTime) isEvent_Type() {}

func (w *Event_PingTime) Unwrap() *PingTime {
	return w.PingTime
}

func (w *Event_PingTime) Pb() dslpingpongpb.Event_Type {
	return &dslpingpongpb.Event_PingTime{PingTime: (w.PingTime).Pb()}
}

func (*Event_PingTime) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.Event_PingTime]()}
}

func EventFromPb(pb *dslpingpongpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
	}
}

func (m *Event) Pb() *dslpingpongpb.Event {
	return &dslpingpongpb.Event{
		Type: (m.Type).Pb(),
	}
}

func (*Event) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.Event]()}
}

type PingTime struct{}

func PingTimeFromPb(pb *dslpingpongpb.PingTime) *PingTime {
	return &PingTime{}
}

func (m *PingTime) Pb() *dslpingpongpb.PingTime {
	return &dslpingpongpb.PingTime{}
}

func (*PingTime) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.PingTime]()}
}

type Message struct {
	Type Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() dslpingpongpb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb dslpingpongpb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *dslpingpongpb.Message_Ping:
		return &Message_Ping{Ping: PingFromPb(pb.Ping)}
	case *dslpingpongpb.Message_Pong:
		return &Message_Pong{Pong: PongFromPb(pb.Pong)}
	}
	return nil
}

type Message_Ping struct {
	Ping *Ping
}

func (*Message_Ping) isMessage_Type() {}

func (w *Message_Ping) Unwrap() *Ping {
	return w.Ping
}

func (w *Message_Ping) Pb() dslpingpongpb.Message_Type {
	return &dslpingpongpb.Message_Ping{Ping: (w.Ping).Pb()}
}

func (*Message_Ping) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.Message_Ping]()}
}

type Message_Pong struct {
	Pong *Pong
}

func (*Message_Pong) isMessage_Type() {}

func (w *Message_Pong) Unwrap() *Pong {
	return w.Pong
}

func (w *Message_Pong) Pb() dslpingpongpb.Message_Type {
	return &dslpingpongpb.Message_Pong{Pong: (w.Pong).Pb()}
}

func (*Message_Pong) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.Message_Pong]()}
}

func MessageFromPb(pb *dslpingpongpb.Message) *Message {
	return &Message{
		Type: Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *dslpingpongpb.Message {
	return &dslpingpongpb.Message{
		Type: (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.Message]()}
}

type Ping struct {
	SeqNr types.SeqNr
}

func PingFromPb(pb *dslpingpongpb.Ping) *Ping {
	return &Ping{
		SeqNr: (types.SeqNr)(pb.SeqNr),
	}
}

func (m *Ping) Pb() *dslpingpongpb.Ping {
	return &dslpingpongpb.Ping{
		SeqNr: (uint64)(m.SeqNr),
	}
}

func (*Ping) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.Ping]()}
}

type Pong struct {
	SeqNr types.SeqNr
}

func PongFromPb(pb *dslpingpongpb.Pong) *Pong {
	return &Pong{
		SeqNr: (types.SeqNr)(pb.SeqNr),
	}
}

func (m *Pong) Pb() *dslpingpongpb.Pong {
	return &dslpingpongpb.Pong{
		SeqNr: (uint64)(m.SeqNr),
	}
}

func (*Pong) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*dslpingpongpb.Pong]()}
}
