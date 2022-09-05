package model

import (
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
)

func PackagePath(sourcePackagePath string) string {
	return sourcePackagePath + "/msgs"
}

func PackageName(sourcePackagePath string) string {
	return sourcePackagePath[strings.LastIndex(sourcePackagePath, "/")+1:] + "msgs"
}

func OutputDir(sourceDir string) string {
	return path.Join(sourceDir, "msgs")
}

// NetMessageNode represents a node in the tree corresponding to the hierarchy of net messages.
type NetMessageNode struct {
	// The protobuf message for this net message.
	message *types.Message
	// The option in the parent's Type oneof.
	oneofOption *types.OneofOption
	// The Type oneof field of the message (if present).
	typeOneof *types.Oneof
	// The children messages in the hierarchy.
	// NB: It may happen that a message class has no children.
	children []*NetMessageNode
	// The parent node in the hierarchy.
	parent *NetMessageNode
	// The accumulated parameters for the constructor function.
	constructorParameters ConstructorParamList
}

// IsRoot returns true if this is the root of the message hierarchy.
func (ev *NetMessageNode) IsRoot() bool {
	return ev.parent == nil && ev.IsMsgClass()
}

// IsMsgClass returns true iff the message has a oneof field marked with [(mir.message_type) = true].
func (ev *NetMessageNode) IsMsgClass() bool {
	return ev.typeOneof != nil
}

// IsNetMessage returns true if this is not a msg class (see IsMsgClass).
func (ev *NetMessageNode) IsNetMessage() bool {
	return !ev.IsMsgClass()
}

// Name returns the name of the message.
// Same as ev.Message().Name().
func (ev *NetMessageNode) Name() string {
	return ev.Message().Name()
}

// Message returns the protobuf message for this net message.
func (ev *NetMessageNode) Message() *types.Message {
	return ev.message
}

// OneofOption returns the option in the parent's Type oneof.
// If nil, IsRoot() must be true.
func (ev *NetMessageNode) OneofOption() *types.OneofOption {
	return ev.oneofOption
}

// TypeOneof returns the Type oneof field of the message (if present).
func (ev *NetMessageNode) TypeOneof() *types.Oneof {
	return ev.typeOneof
}

// Children returns the children messages in the hierarchy.
// NB: It may happen that a message class has no children.
func (ev *NetMessageNode) Children() []*NetMessageNode {
	return ev.children
}

// Parent returns the parent event in the hierarchy.
func (ev *NetMessageNode) Parent() *NetMessageNode {
	return ev.parent
}

// AllConstructorParameters returns the accumulated parameters for the constructor function.
// The parameters include all the fields of all the ancestors in the hierarchy except those marked with
// [(mir.omit_in_constructor) = true] and the Type oneofs.
// To get the parameters that correspond to the fields only of this node
func (ev *NetMessageNode) AllConstructorParameters() ConstructorParamList {
	return ev.constructorParameters
}

// ThisNodeConstructorParameters returns a suffix of AllConstructorParameters() that corresponds to the fields
// only of this in the hierarchy, without the fields accumulated from the ancestors.
func (ev *NetMessageNode) ThisNodeConstructorParameters() ConstructorParamList {
	if ev.Parent() == nil {
		return ev.AllConstructorParameters()
	}

	// Remove the prefix that corresponds to the parameters of the parent.
	return ConstructorParamList{ev.constructorParameters.Slice[len(ev.Parent().constructorParameters.Slice):]}
}

func (ev *NetMessageNode) Constructor() *jen.Statement {
	return jen.Qual(PackagePath(ev.Message().PbPkgPath()), ev.Name())
}
