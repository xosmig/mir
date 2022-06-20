package contextstore

import "fmt"

type sequentialContextStoreImpl[T any] struct {
	nextItemID ItemID
	storage    map[ItemID]T
	id         uint64
	hasID      bool
}

// NewSequentialContextStore creates an empty ContextStore that can only be accessed sequentially.
func NewSequentialContextStore[T any]() ContextStore[T] {
	return &sequentialContextStoreImpl[T]{
		storage: make(map[ItemID]T),
	}
}

// Store stores the given data in the ContextStore and returns a unique id.
// The data can be later recovered or disposed of using this id.
func (s *sequentialContextStoreImpl[T]) Store(t T) ItemID {
	id := s.nextItemID
	s.nextItemID++

	s.storage[id] = t
	return id
}

// Recover returns the data stored under the provided id.
// Note that the data will continue to exist in the ContextStore.
// In order to dispose of the data, call s.Dispose(id) or s.RecoverAndDispose(id).
func (s *sequentialContextStoreImpl[T]) Recover(id ItemID) T {
	item, present := s.storage[id]
	if !present {
		panic(fmt.Errorf("item with id '%v' is not present in the ContextStore", id))
	}

	return item
}

// Dispose removes the data from the ContextStore.
func (s *sequentialContextStoreImpl[T]) Dispose(id ItemID) {
	delete(s.storage, id)
}

// RecoverAndDispose returns the data stored under the provided id and removes it from the ContextStore.
func (s *sequentialContextStoreImpl[T]) RecoverAndDispose(id ItemID) T {
	t := s.Recover(id)
	s.Dispose(id)
	return t
}

// SetID assigns an identifier to the ContextStore.
// IDs can be useful to manage multiple context stores in a single module.
func (s *sequentialContextStoreImpl[T]) SetID(id uint64) {
	s.id = id
	s.hasID = true
}

// HasID returns true iff SetID was previously invoked.
func (s *sequentialContextStoreImpl[T]) HasID() bool {
	return s.hasID
}

// GetID returns the id previously assigned by SetID.
// If SetID was not called before GetID, GetID panics.
func (s *sequentialContextStoreImpl[T]) GetID() uint64 {
	return s.id
}
