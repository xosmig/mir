package availabilitypb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_RequestBatch) Unwrap() *RequestBatch {
	return p.RequestBatch
}

func (p *Event_NewBatch) Unwrap() *NewBatch {
	return p.NewBatch
}
