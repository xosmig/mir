package dslpingpongpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[T any] interface {
	Event_Type
	Unwrap() *T
}

func (w *Event_PingTime) Unwrap() *PingTime {
	return w.PingTime
}

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_Ping) Unwrap() *Ping {
	return w.Ping
}

func (w *Message_Pong) Unwrap() *Pong {
	return w.Pong
}
