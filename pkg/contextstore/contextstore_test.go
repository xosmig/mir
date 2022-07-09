package contextstore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequentialContextStoreImpl_RecoverAndDispose(t *testing.T) {
	cs := NewSequentialContextStore[string]()
	helloID := cs.Store("Hello")
	worldID := cs.Store("World")

	item, ok := cs.Recover(worldID)
	assert.True(t, ok)
	assert.Equal(t, "World", item)

	item, ok = cs.Recover(helloID)
	assert.True(t, ok)
	assert.Equal(t, "Hello", item)

	cs.Dispose(worldID)

	item, ok = cs.Recover(worldID)
	assert.False(t, ok)
	assert.Equal(t, "", item)

	item, ok = cs.RecoverAndDispose(helloID)
	assert.True(t, ok)
	assert.Equal(t, "Hello", item)

	item, ok = cs.RecoverAndDispose(helloID)
	assert.False(t, ok)
	assert.Equal(t, "", item)

	assert.NotPanics(t, func() {
		cs.Dispose(worldID)
	})

	assert.NotPanics(t, func() {
		cs.Dispose(helloID)
	})
}
