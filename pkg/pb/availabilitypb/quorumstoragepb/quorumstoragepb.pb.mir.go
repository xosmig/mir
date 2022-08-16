package quorumstoragepb

type Event_Type = isEvent_Type

type Event_TypeWrapper[Ev any] interface {
	Event_Type
	Unwrap() *Ev
}

func (p *Event_Store) Unwrap() *StoreBatchOnQuorum {
	return p.Store
}

func (p *Event_Stored) Unwrap() *BatchStoredOnQuorum {
	return p.Stored
}

func (p *Event_VerifyBatch) Unwrap() *VerifyBatch {
	return p.VerifyBatch
}

func (p *Event_BatchVerified) Unwrap() *BatchVerified {
	return p.BatchVerified
}

func (p *Event_VerifyCert) Unwrap() *VerifyCert {
	return p.VerifyCert
}

func (p *Event_CertVerified) Unwrap() *CertVerified {
	return p.CertVerified
}

func (p *Event_RetrieveTransactions) Unwrap() *RetrieveTransactions {
	return p.RetrieveTransactions
}

func (p *Event_TransactionsRetrieved) Unwrap() *TransactionsRetrieved {
	return p.TransactionsRetrieved
}
