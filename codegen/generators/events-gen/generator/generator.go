package generator

import (
	"fmt"
	"reflect"

	"github.com/filecoin-project/mir/codegen/generators/events-gen/events"

	"github.com/filecoin-project/mir/codegen/generators/types-gen/types"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
)

type EventsGenerator struct{}

func (EventsGenerator) Run(pbGoStructTypes []reflect.Type) error {
	eventRootMessages, err := GetEventRootMessagesUsingDefaultParser(pbGoStructTypes)
	if err != nil {
		return err
	}

	for _, eventRootMessage := range eventRootMessages {
		eventParser := events.DefaultParser()

		eventRoot, err := eventParser.ParseEventHierarchy(eventRootMessage)
		if err != nil {
			return err
		}

		err = GenerateEventConstructors(eventRoot)
		if err != nil {
			return fmt.Errorf("error generating event constructors: %w", err)
		}

		//err = GenerateDslFunctions(eventRoot)
		//if err != nil {
		//	return fmt.Errorf("error generating dsl functions: %w", err)
		//}
	}

	return nil
}

func GetEventRootMessagesUsingDefaultParser(pbGoStructTypes []reflect.Type) ([]*types.Message, error) {
	typesParser := types.DefaultParser()

	// For convenience, the parser operates on pointer to struct types and not struct types themselves.
	// The reason for this is that protobuf messages are always used as pointers in Go code.
	ptrTypes := sliceutil.Transform(pbGoStructTypes, func(_ int, tp reflect.Type) reflect.Type {
		return reflect.PointerTo(tp)
	})

	msgs, err := typesParser.ParseMessages(ptrTypes)
	if err != nil {
		return nil, err
	}

	// Look for the root of the event hierarchy.
	return sliceutil.Filter(msgs, func(_ int, msg *types.Message) bool { return msg.IsEventRoot() }), nil
}
