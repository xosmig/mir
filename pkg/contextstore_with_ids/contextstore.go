package contextstore_with_ids

// ContextStore can be used to store arbitrary data under an automatically deterministically generated unique id.
type ContextStore[T any] interface {
	Id() StoreID
	Store(t T) ItemID
	Recover(id ItemID) T
	Dispose(id ItemID)
	RecoverAndDispose(id ItemID) T
}

// StoreID is used to uniquely identify a ContextStore when several of them are used in a single module.
type StoreID uint64

// Pb returns the protobuf representation of StoreID.
func (i StoreID) Pb() uint64 {
	return uint64(i)
}

// ItemID is used to uniquely identify entries of the ContextStore.
type ItemID uint64

// Pb returns the protobuf representation of ItemID.
func (i ItemID) Pb() uint64 {
	return uint64(i)
}
