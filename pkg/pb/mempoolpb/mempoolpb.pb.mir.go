package mempoolpb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_NewTransactions) Unwrap() *NewTransactions {
	return p.NewTransactions
}

func (p *Event_RequestTransactions) Unwrap() *RequestTransactions {
	return p.RequestTransactions
}

func (p *Event_TransactionsResponse) Unwrap() *TransactionsResponse {
	return p.TransactionsResponse
}
