package contextstore_with_ids

import "fmt"

type sequentialContextStoreImpl[T any] struct {
	id         StoreID
	nextItemID ItemID
	storage    map[ItemID]T
}

// NewSequentialContextStore creates an empty ContextStore that can only be accessed sequentially.
func NewSequentialContextStore[T any](storeID uint64) ContextStore[T] {
	return &sequentialContextStoreImpl[T]{
		id:         StoreID(storeID),
		nextItemID: 0,
		storage:    make(map[ItemID]T),
	}
}

func (s *sequentialContextStoreImpl[T]) Id() StoreID {
	return s.id
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
