package dsl

import (
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/util/reflectutil"
)

// The functions in this file provide low-level access to the internals of a dsl module.
// The should not be used by most applications. However, they can be useful in order to extend the
// functionality of dsl modules.

// GetEventHandlers returns the slice of registered event handlers for event
func GetEventHandlers[EvWrapper eventpb.Event_Type](h Handle) []EventHandler {
	return h.impl.eventHandlers[reflectutil.TypeOf[EvWrapper]()]
}
