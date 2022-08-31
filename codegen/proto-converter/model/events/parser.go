package events

import (
	"fmt"

	"github.com/filecoin-project/mir/codegen/proto-converter/model/types"
)

type Parser struct {
	messageParser *types.Parser
}

func NewParser(messageParser *types.Parser) *Parser {
	return &Parser{messageParser}
}

// ParseEventHierarchy extracts the information about the whole event hierarchy by its root.
func (p *Parser) ParseEventHierarchy(eventRootMsg *types.Message) (root *EventNode, err error) {
	if !eventRootMsg.IsEventRoot() {
		return nil, fmt.Errorf("message %v is not marked as event root", eventRootMsg.Name())
	}

	root, err = p.parseEventNodeRecursively(eventRootMsg, nil, nil, ConstructorParamList{})
	return
}

// parseEventNodeRecursively parses a message from the event hierarchy.
// parent is the parent in the event hierarchy. Note that the parent's list of children may not be complete.
func (p *Parser) parseEventNodeRecursively(
	msg *types.Message,
	optionInParentOneof *types.OneofOption,
	parent *EventNode,
	constructorParameters ConstructorParamList,
) (*EventNode, error) {

	fields, err := p.messageParser.ParseFields(msg)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		if field.IsEventTypeOneof() || field.OmitInConstructor() {
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
				return nil, fmt.Errorf("non-message type in the event hierarchy: %v", opt.Field.Name)
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
