package eventutil

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"reflect"
)

// GetOrigin uses reflection to descend the hierarchy of events by using
// repeatedly the field named "Type" and then returns the value of the field
// named "Origin".
//
// This function exists solely to save some boilerplate in unit tests and make
// them more readable.
func GetOrigin[Origin any](ev *eventpb.Event) *Origin {
	// cur is a reflect.Value of a pointer to a protobuf-generated struct.
	cur := reflect.ValueOf(ev)
	for {
		if cur.Kind() != reflect.Pointer || cur.Elem().Kind() != reflect.Struct {
			panic(fmt.Sprintf("expected a pointer to a struct, got %T", cur.Interface()))
		}
		typeField, ok := cur.Elem().Type().FieldByName("Type")
		if !ok {
			break
		}

		protoOneofWrapper := reflect.ValueOf(cur.Elem().FieldByIndex(typeField.Index).Interface())

		unwrapMethod, ok := protoOneofWrapper.Type().MethodByName("Unwrap")
		if !ok {
			panic(fmt.Sprintf("oneof Type is not marked as mir.event_type and has no Unwrap() method"))
		}

		cur = protoOneofWrapper.Method(unwrapMethod.Index).Call([]reflect.Value{})[0]
	}

	if _, ok := cur.Elem().Type().FieldByName("Origin"); !ok {
		panic(fmt.Sprintf("event %v does not have a field named 'Origin'", ev))
	}

	origin := cur.Elem().FieldByName("Origin").Interface()
	res, ok := origin.(*Origin)
	if !ok {
		panic(fmt.Sprintf("event %v has field 'Origin', but it has type %T and not %T", ev, origin, res))
	}

	return res
}
