package model

import "github.com/filecoin-project/mir/codegen/proto-converter/util/jenutil"

// TypeOneofFieldName is the name of the oneof field that lists the children of the node in the event hierarchy.
const TypeOneofFieldName = "Type"

type EventNode struct {
	// The message for this event.
	message *Message
	// The option in the parent's Type oneof.
	oneofOption *OneofOption
	// The Type oneof field of the message (if present).
	typeOneof *Oneof
	// The children events in the hierarchy.
	// NB: It may happen that an event class has no children.
	children []*EventNode
	// The parent event in the hierarchy.
	parent *EventNode
	// The accumulated parameters for the constructor function.
	constructorParameters *jenutil.FuncParamList
}

// IsRoot returns true if this is the root of the event hierarchy.
func (ev *EventNode) IsRoot() bool {
	return ev.parent == nil && ev.IsEventClass()
}

// IsEventClass returns true iff the message has a oneof field named Type.
func (ev *EventNode) IsEventClass() bool {
	return ev.typeOneof != nil
}

// IsEvent returns true if this is an event (i.e., the message has no oneof field named Type).
func (ev *EventNode) IsEvent() bool {
	return !ev.IsEventClass()
}

// Message returns the message for this event.
func (ev *EventNode) Message() *Message {
	return ev.message
}

// OneofOption returns the option in the parent's Type oneof.
// If nil, IsRoot() must be true.
func (ev *EventNode) OneofOption() *OneofOption {
	return ev.oneofOption
}

// TypeOneof returns the Type oneof field of the message (if present).
func (ev *EventNode) TypeOneof() *Oneof {
	return ev.typeOneof
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

// ConstructorParameters returns the accumulated parameters for the constructor function.
func (ev *EventNode) ConstructorParameters() *jenutil.FuncParamList {
	return ev.constructorParameters
}
