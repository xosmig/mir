package bcbpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_Request)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
	}
}
