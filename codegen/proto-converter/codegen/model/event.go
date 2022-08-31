package model

type EventNode struct {
	// The message for this event.
	message *Message
	// Whether the message has a oneof field named Type.
	isEventClass bool
	// The children events in the hierarchy.
	// NB: It may happen that an event class has no children.
	children []*EventNode
	// The parent event in the hierarchy.
	parent *EventNode
}

// IsRoot returns true if this is the root of the event hierarchy.
// IsRoot() implies IsEventClass().
func (ev *EventNode) IsRoot() bool {
	return ev.parent == nil && ev.isEventClass
}

// IsEventClass returns true iff the message has a oneof field named Type.
func (ev *EventNode) IsEventClass() bool {
	return ev.isEventClass
}

// IsEvent returns true if this is an event (i.e., the message has no oneof field named Type).
func (ev *EventNode) IsEvent() bool {
	return !ev.isEventClass
}

// Message returns the message for this event.
func (ev *EventNode) Message() *Message {
	return ev.message
}

// Children returns the children events in the hierarchy.
// NB: It may happen that an event class has no children.
func (ev *EventNode) Children() []*EventNode {
	return ev.children
}

// Parent returns the parent event in the hierarchy.
func (ev *EventNode) Parent() *EventNode {
	return ev.parent
}
