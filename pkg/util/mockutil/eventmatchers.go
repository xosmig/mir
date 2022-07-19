package mockutil

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/util/reflectutil"
	"github.com/golang/mock/gomock"
	"reflect"
)

// EventOfType matches events with the given event type.
// exaxmple:
//     testModule.EXPECT().Event(mockutil.EventOfType[*eventpb.Event_TestingString]())
func EventOfType[EvType eventpb.Event_Type]() gomock.Matcher {
	return &eventTypeMatcher[EvType]{}
}

type eventTypeMatcher[EvType eventpb.Event_Type] struct{}

// Matches returns whether x is a match.
func (m *eventTypeMatcher[EvType]) Matches(x any) bool {
	ev, ok := x.(*eventpb.Event)
	if !ok {
		return false
	}

	_, ok = ((any)(ev.Type)).(*EvType)
	return true
}

// String describes what the matcher matches.
func (m *eventTypeMatcher[EvType]) String() string {
	return fmt.Sprintf("{Event of type %v}", reflectutil.TypeOf[EvType]().Name())
}

// EventOfSubtype matches events with the given event type and the given subtype.
// The subtype is determined using the field named Type of the event.
// Note that the third generic parameter (Ev) can always be inferred automatically from the first (EvType)
// (see the example below).
// example:
//     mempool.EXPECT().Event(mockutil.EventOfSubtype[*eventpb.Event_Mempool, *mempoolpb.Event_RequestBatch]())
func EventOfSubtype[EvType eventpb.Event_TypeWrapper[Ev], EvSubtype any, Ev any]() gomock.Matcher {
	typeTp := reflectutil.TypeOf[EvType]()
	subtypeTp := reflectutil.TypeOf[EvSubtype]()

	field, ok := reflectutil.TypeOf[Ev]().FieldByName("Type")
	if !ok {
		panic(fmt.Errorf("event type %v does not have subtypes", reflectutil.TypeOf[Ev]().Name()))
	}

	if !subtypeTp.Implements(field.Type) {
		panic(fmt.Errorf("%v is not a subtype of %v", subtypeTp.Name(), typeTp.Name()))
	}

	return &eventSubtypeMatcher[EvType, Ev]{field, subtypeTp}
}

type eventSubtypeMatcher[EvType eventpb.Event_TypeWrapper[Ev], Ev any] struct {
	field     reflect.StructField
	subtypeTp reflect.Type
}

// Matches returns whether x is a match.
func (m *eventSubtypeMatcher[EvType, Ev]) Matches(x any) bool {
	event, ok := x.(*eventpb.Event)
	if !ok {
		return false
	}

	evWrapper, ok := ((any)(event.Type)).(*EvType)
	ev := ((any)(evWrapper)).(eventpb.Event_TypeWrapper[Ev]).Unwrap()

	return reflect.ValueOf(ev).FieldByName("Type").CanConvert(m.subtypeTp)
}

// String describes what the matcher matches.
func (m *eventSubtypeMatcher[EvType, Ev]) String() string {
	return fmt.Sprintf("{Event of type %v of subtype %v}", reflectutil.TypeOf[EvType]().Name(), m.subtypeTp.Name())
}
