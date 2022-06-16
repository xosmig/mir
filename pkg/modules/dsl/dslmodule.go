package dsl

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"reflect"
	"unsafe"
)

// dslModuleImpl allows creating passive modules in a very natural declarative way.
type dslModuleImpl struct {
	eventHandlers     map[reflect.Type][]func(ev *eventpb.Event) error
	conditionHandlers []func() error
	outputEvents      *events.EventList
}

type Handle struct {
	moduleImpl *dslModuleImpl
}

type DslModule interface {
	GetDslHandle() Handle
	GetPassiveModule() modules.PassiveModule
}

func (m *dslModuleImpl) GetDslHandle() Handle {
	return Handle{moduleImpl: m}
}

func (m *dslModuleImpl) GetPassiveModule() modules.PassiveModule {
	return m
}

func NewModule() *dslModuleImpl {
	return &dslModuleImpl{}
}

//type conditionEntry struct {
//	condition func() bool
//	handler   func() error
//}

type evContainer[Ev any] struct{ ev *Ev }

// typeOf returns the reflect.Type that represents the type T.
// TODO: consider moving this function to some utility package. Ideally, it should be in the standard library :)
// see: https://github.com/golang/go/issues/50741
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// UponEvent registers an event handler for module m.
// This event handler will be called every time an event of type EvTp is received.
// Type EvTp is the protoc-generated wrapper around Ev -- protobuf representation of the event.
// Note that the type parameter Ev can be inferred automatically from handler.
func UponEvent[EvTp, Ev any](m DslModule, handler func(ev *Ev) error) {
	evTpType := typeOf[EvTp]()
	evType := typeOf[Ev]()
	evContainerType := typeOf[evContainer[Ev]]()

	m.GetDslHandle().moduleImpl.eventHandlers[evTpType] = append(
		m.GetDslHandle().moduleImpl.eventHandlers[evTpType],
		func(ev *eventpb.Event) error {
			evTp := ev.Type.(EvTp)
			// The safety of this cast is verified by the runtime checks below.
			evTpPtr := (*evContainer[Ev])(unsafe.Pointer(&evTp))
			return handler((*evTpPtr).ev)
		})

	// These checks verify that an object of type EvTp can be safely interpreted as object of type evContainer[Ev].
	// They are only performed at the time of registration of the handler (which is supposedly at the very beginning
	// of the program execution) and make sure that no unexpected runtime errors will happen during later stages
	// of protocol execution and no memory will be corrupted by the pointer cast between (*EvTp) and (*evContainer[Ev]).
	// These checks may fail in case of a change to the protobuf API or when erroneous type parameters (EvTp and Ev) are
	// provided.
	const explanationStr = "Most likely, the function was called with an invalid combination of type parameters. " +
		"This error may also be caused by a change to the API of the protobuf implementation."
	if evTpType.Kind() != reflect.Struct || evTpType.NumField() != 1 {
		panic(fmt.Sprintf("%s is supposed to be a struct with a single field of type *%s. %s",
			evTpType.Name(), evType.Name(), explanationStr))
	}
	if evTpType.Field(0).Type.Kind() != reflect.Pointer || evTpType.Field(0).Type.Elem() != evType {
		panic(fmt.Sprintf("%s is supposed to be a struct with a single field of type *%s. %s",
			evTpType.Name(), evType.Name(), explanationStr))
	}
	if evTpType.Field(0).Offset != evContainerType.Field(0).Offset {
		panic(fmt.Sprintf("Unexpected field offset for type %s. %s",
			evTpType.Name(), explanationStr))
	}
	if evTpType.Size() != evContainerType.Size() || evTpType.Align() != evContainerType.Align() {
		panic(fmt.Sprintf("Unexpected size or alignment for type %s. %s",
			evTpType.Name(), explanationStr))
	}
}

//// UponEventWithCondition registers a conditional event handler for module m.
//// The condition cond will be checked every time an event of type EvTp is received and the handler will be invoked
//// whenever the condition evaluates to true.
//// Type EvTp is the protoc-generated wrapper around Ev -- protobuf representation of the event.
//// Note that the type parameter Ev can be inferred automatically from handler.
//// TODO: decide whether we want to keep this function.
//func UponEventWithCondition[EvTp, Ev any](m DslModule, cond func(ev *Ev) bool, handler func(ev *Ev) error) {
//	UponEvent[EvTp](m, func(ev *Ev) error {
//		if !cond(ev) {
//			return nil
//		}
//		return handler(ev)
//	})
//}

//// UponRepeatedCondition registers a *repeated* condition handler. Predicate `cond` will be evaluated each time after a
//// batch of events is processed and *each time* `cond()` returns `True`, `handler` will be invoked.
//// Conditions are checked in the order of their registration.
//func UponRepeatedCondition(m DslModule, cond func() bool, handler func() error) {
//	m.GetDslHandle().moduleImpl.conditionHandlers = append(
//		m.GetDslHandle().moduleImpl.conditionHandlers,
//		conditionEntry{
//			condition: cond,
//			handler:   handler,
//		})
//}

//// UponOneShotCondition registers a *one-shot* condition handler. Predicate `cond` will be evaluated each time after a
//// batch of events is processed until it returns `true`. After `cond()` returns `true` *for the first time*  `cond()`,
//// `handler` will be invoked.
//// Conditions are checked in the order of their registration.
//func UponOneShotCondition(m DslModule, cond func() bool, handler func() error) {
//	// Note that a more efficient implementation that actually removes the one-shot condition from the list is possible.
//	fired := false
//	UponRepeatedCondition(m, func() bool { return !fired && cond() }, func() error {
//		fired = true
//		return handler()
//	})
//}

// UponCondition registers a special type of handler that will be invoked each time after processing a batch of events.
// The handler is assumed to represent a conditional action: it is supposed to check some predicate on the state
// and perform actions if the predicate evaluates is satisfied.
func UponCondition(m DslModule, handler func() error) {
	m.GetDslHandle().moduleImpl.conditionHandlers = append(m.GetDslHandle().moduleImpl.conditionHandlers, handler)
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (m *dslModuleImpl) ImplementsModule() {}

func EmitEvent(m DslModule, ev *eventpb.Event) {
	m.GetDslHandle().moduleImpl.outputEvents.PushBack(ev)
}

func (m *dslModuleImpl) ApplyEvents(evs *events.EventList) (*events.EventList, error) {
	// Run event handlers.
	iter := evs.Iterator()
	for ev := iter.Next(); ev != nil; ev = iter.Next() {
		handlers, ok := m.eventHandlers[reflect.TypeOf(ev.Type)]
		if !ok {
			return nil, fmt.Errorf("unknown event type '%T'", ev.Type)
		}

		for _, h := range handlers {
			err := h(ev)
			if err != nil {
				return nil, err
			}
		}
	}

	//// Run condition handlers.
	//for _, condEntry := range m.conditionHandlers {
	//	if condEntry.condition() {
	//		err := condEntry.handler()
	//
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	//}

	// Run condition handlers.
	for _, condition := range m.conditionHandlers {
		err := condition()

		if err != nil {
			return nil, err
		}
	}

	outputEvents := m.outputEvents
	m.outputEvents = &events.EventList{}
	return outputEvents, nil
}
