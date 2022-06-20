package contextstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequentialContextStoreImpl_RecoverAndDispose(t *testing.T) {
	cs := NewSequentialContextStore[string]()
	helloID := cs.Store("Hello")
	worldID := cs.Store("World")

	assert.Equal(t, "World", cs.Recover(worldID))
	assert.Equal(t, "Hello", cs.Recover(helloID))

	cs.Dispose(worldID)
	assert.Panics(t, func() {
		cs.Recover(worldID)
	})

	assert.Equal(t, "Hello", cs.RecoverAndDispose(helloID))
	assert.Panics(t, func() {
		cs.RecoverAndDispose(helloID)
	})
}

func TestSequentialContextStoreImpl_GetID(t *testing.T) {
	cs := NewSequentialContextStore[string]()

	assert.Panics(t, func() {
		cs.GetID()
	})
	cs.SetID(17)
	assert.Equal(t, 17, cs.GetID())
}
