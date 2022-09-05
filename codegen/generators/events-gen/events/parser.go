package events

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/filecoin-project/mir/codegen/generators/types-gen/types"
	"github.com/filecoin-project/mir/pkg/pb/mir/dsl"
)

type Parser struct {
	messageParser  *types.Parser
	eventNodeCache map[reflect.Type]*EventNode
}

var defaultParser = newParser(types.DefaultParser())

// DefaultParser returns a singleton Parser.
// It must not be accessed concurrently.
func DefaultParser() *Parser {
	return defaultParser
}

// newParser is not exported as DefaultParser() is supposed to be used instead.
func newParser(messageParser *types.Parser) *Parser {
	return &Parser{
		messageParser:  messageParser,
		eventNodeCache: make(map[reflect.Type]*EventNode),
	}
}

// ParseEventHierarchy extracts the information about the whole event hierarchy by its root.
func (p *Parser) ParseEventHierarchy(eventRootMsg *types.Message) (root *EventNode, err error) {

	if !eventRootMsg.IsEventRoot() {
		return nil, fmt.Errorf("message %v is not marked as event root", eventRootMsg.Name())
	}

	root, err = p.parseEventNodeRecursively(eventRootMsg, nil, nil, types.ConstructorParamList{})
	return
}

// parseEventNodeRecursively parses a message from the event hierarchy.
// parent is the parent in the event hierarchy. Note that the parent's list of children may not be complete.
func (p *Parser) parseEventNodeRecursively(
	msg *types.Message,
	optionInParentOneof *types.OneofOption,
	parent *EventNode,
	constructorParameters types.ConstructorParamList,
) (node *EventNode, err error) {

	// First, check the cache.
	if tp, ok := p.eventNodeCache[msg.PbReflectType()]; ok {
		return tp, nil
	}

	// Remember the result in the cache when finished
	defer func() {
		if err == nil && node != nil {
			p.eventNodeCache[msg.PbReflectType()] = node
		}
	}()

	fields, err := p.messageParser.ParseFields(msg)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		if field.IsEventTypeOneof() || DslIgnore(field.ProtoDesc) {
			continue
		}

		constructorParameters = constructorParameters.Append(field.LowercaseName(), field)
	}

	// Check if this is an event class.
	if typeOneof, ok := getTypeOneof(fields); ok {
		node := &EventNode{
			message:               msg,
			oneofOption:           optionInParentOneof,
			typeOneof:             typeOneof,
			children:              nil, // to be filled separately
			parent:                parent,
			constructorParameters: constructorParameters,
		}

		for _, opt := range typeOneof.Options {
			childMsg, ok := opt.Field.Type.(*types.Message)
			if !ok {
				return nil, fmt.Errorf("non-message type in the event hierarchy: %v", opt.Name())
			}

			if !childMsg.ShouldGenerateMirType() {
				// Skip children that are not marked as Mir
				continue
			}

			childNode, err := p.parseEventNodeRecursively(childMsg, opt, node, constructorParameters)
			if err != nil {
				return nil, err
			}

			node.children = append(node.children, childNode)
		}

		return node, nil
	}

	return &EventNode{
		message:               msg,
		oneofOption:           optionInParentOneof,
		typeOneof:             nil,
		children:              nil,
		parent:                parent,
		constructorParameters: constructorParameters,
	}, nil
}

func getTypeOneof(fields types.Fields) (*types.Oneof, bool) {
	for _, field := range fields {
		// Recursively call the generator on all subtypes.
		if field.IsEventTypeOneof() {
			return field.Type.(*types.Oneof), true
		}
	}
	return nil, false
}

func DslIgnore(protoDesc protoreflect.Descriptor) bool {
	fieldDesc, ok := protoDesc.(protoreflect.FieldDescriptor)
	if !ok {
		// This is a oneof and this option does not exist for oneofs.
		// Return false as the default value.
		return false
	}

	return proto.GetExtension(fieldDesc.Options().(*descriptorpb.FieldOptions), dsl.E_Ignore).(bool)
}
