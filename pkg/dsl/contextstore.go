package dsl

import (
	cs "github.com/filecoin-project/mir/pkg/contextstore"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
)

type ContextStore[T any] struct {
	impl    cs.ContextStore[T]
	storeID uint64
	hasID   bool
}

func NewDSLContextStore[T any]() *ContextStore[T] {
	return &ContextStore[T]{
		impl: cs.NewSequentialContextStore[T](),
	}
}

func (m *dslModuleImpl) getOrSetContextStoreId(contextStore interface{}) uint64 {
	id, ok := m.contextStoreIDs[contextStore]
	if !ok {
		id = m.nextStoreID
		m.nextStoreID++
		m.contextStoreIDs[id] = m.nextStoreID
	}

	return id
}

func dslOrigin(contextStoreID uint64, itemID cs.ItemID) *eventpb.DslOrigin {
	return &eventpb.DslOrigin{
		ContextStoreId: contextStoreID,
		ItemID:         itemID.Pb(),
	}
}
