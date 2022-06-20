package dsl

import (
	cs "github.com/filecoin-project/mir/pkg/contextstore"
)

// This type is not exported in order to prevent the users from instantiating it directly.
// The only valid way to instantiate a ContextHandle is by invoking NewContextHandle(m).
type contextHandleImpl[C any] struct {
	store     cs.ContextStore[C]
	id        ContextHandleID
	dslModule Module
}

type ContextHandleID uint64

func (id ContextHandleID) Pb() uint64 {
	return uint64(id)
}

type ContextHandle[C any] *contextHandleImpl[C]

func (h ContextHandle[C]) StoreContext(context C) cs.ItemID {
	return h.store.Store(context)
}

func (h ContextHandle[C]) RecoverContext(itemID cs.ItemID) C {
	return h.store.Recover(itemID)
}

func (h ContextHandle[C]) DeferCleanup(itemID cs.ItemID) {
	dslHandle := h.dslModule.GetDslHandle()
	dslHandle.eventCleanup = append(dslHandle.eventCleanup, func() {
		h.store.Dispose(itemID)
	})
}

func (h ContextHandle[C]) RecoverContextAndDeferCleanup(itemID cs.ItemID) C {
	res := h.RecoverContext(itemID)
	h.DeferCleanup(itemID)
	return res
}

func (h ContextHandle[C]) GetID() ContextHandleID {
	return h.id
}

func NewContextHandle[C any](m Module) ContextHandle[C] {
	id := m.GetDslHandle().nextContextHandleID
	m.GetDslHandle().nextContextHandleID++

	return &contextHandleImpl[C]{
		store: cs.NewSequentialContextStore[C](),
		id:    ContextHandleID(id),
	}
}
