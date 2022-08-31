package events

import (
	"path"
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/filecoin-project/mir/codegen/proto-converter/model/types"
)

func PackagePath(sourcePackagePath string) string {
	return sourcePackagePath + "/events"
}

func PackageName(sourcePackagePath string) string {
	return sourcePackagePath[strings.LastIndex(sourcePackagePath, "/")+1:] + "events"
}

func OutputDir(sourceDir string) string {
	return path.Join(sourceDir, "events")
}

type EventNode struct {
	// The message for this event.
	message *types.Message
	// The option in the parent's Type oneof.
	oneofOption *types.OneofOption
	// The Type oneof field of the message (if present).
	typeOneof *types.Oneof
	// The children events in the hierarchy.
	// NB: It may happen that an event class has no children.
	children []*EventNode
	// The parent event in the hierarchy.
	parent *EventNode
	// The accumulated parameters for the constructor function.
	constructorParameters ConstructorParamList
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

// Name returns the name of the event, which is the name of the corresponding message.
func (ev *EventNode) Name() string {
	return ev.Message().Name()
}

// Message returns the message for this event.
func (ev *EventNode) Message() *types.Message {
	return ev.message
}

// OneofOption returns the option in the parent's Type oneof.
// If nil, IsRoot() must be true.
func (ev *EventNode) OneofOption() *types.OneofOption {
	return ev.oneofOption
}

// TypeOneof returns the Type oneof field of the message (if present).
func (ev *EventNode) TypeOneof() *types.Oneof {
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

// AllConstructorParameters returns the accumulated parameters for the constructor function.
// The parameters include all the fields of all the ancestors in the hierarchy except those marked with
// [(mir.omit_in_constructor) = true] and the Type oneofs.
// To get the parameters that correspond to the fields only of this node
func (ev *EventNode) AllConstructorParameters() ConstructorParamList {
	return ev.constructorParameters
}

// ThisNodeConstructorParameters returns a suffix of AllConstructorParameters() that corresponds to the fields
// only of this in the hierarchy, without the fields accumulated from the ancestors.
func (ev *EventNode) ThisNodeConstructorParameters() ConstructorParamList {
	if ev.Parent() == nil {
		return ev.AllConstructorParameters()
	}

	// Remove the prefix that corresponds to the parameters of the parent.
	return ConstructorParamList{ev.constructorParameters.Slice[len(ev.Parent().constructorParameters.Slice):]}
}

func (ev *EventNode) Constructor() *jen.Statement {
	return jen.Qual(PackagePath(ev.Message().PbPkgPath()), ev.Name())
}
