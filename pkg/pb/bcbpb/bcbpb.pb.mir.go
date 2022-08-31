package bcbpb

import (
	reflect "reflect"
)

type Event_Type = isEvent_Type

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_Request)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
	}
}
