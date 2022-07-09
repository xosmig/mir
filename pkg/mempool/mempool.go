package mempool

// TxID is a unique identifier of a transaction.
type TxID string

func (id TxID) Pb() string {
	return string(id)
}

// BatchID is a unique identifier of a batch.
type BatchID string

func (id BatchID) Pb() string {
	return string(id)
}
