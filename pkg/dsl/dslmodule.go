package dsl

import (
	"fmt"
	cs "github.com/filecoin-project/mir/pkg/contextstore"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/reflectutil"
	"reflect"
	"unsafe"
)

// dslModuleImpl allows creating passive modules in a very natural declarative way.
type dslModuleImpl struct {
	moduleID          t.ModuleID
	eventHandlers     map[reflect.Type][]func(ev *eventpb.Event) error
	conditionHandlers []func() error
	outputEvents      *events.EventList
	// contextStore is used to store and recover context on asynchronous operations such as signature verification.
	contextStore cs.ContextStore[any]
}

type Handle struct {
	impl *dslModuleImpl
}

type Module interface {
	modules.PassiveModule
	GetDslHandle() Handle
}

func (m *dslModuleImpl) GetDslHandle() Handle {
	return Handle{m}
}

func NewModule(moduleID t.ModuleID) *dslModuleImpl {
	return &dslModuleImpl{
		moduleID:      moduleID,
		eventHandlers: make(map[reflect.Type][]func(ev *eventpb.Event) error),
		outputEvents:  &events.EventList{},
		contextStore:  cs.NewSequentialContextStore[any](),
	}
}

type evContainer[Ev any] struct{ ev *Ev }

// UponEvent registers an event handler for module m.
// This event handler will be called every time an event of type Ev is received.
// Type EvWrapper is the protoc-generated wrapper around Ev -- protobuf representation of the event.
// Note that the type parameter Ev can be inferred automatically from handler.
func UponEvent[EvWrapper, Ev any](m Module, handler func(ev *Ev) error) {
	evTpType := reflectutil.TypeOf[EvWrapper]()
	evType := reflectutil.TypeOf[Ev]()
	evContainerType := reflectutil.TypeOf[evContainer[Ev]]()

	m.GetDslHandle().impl.eventHandlers[evTpType] = append(
		m.GetDslHandle().impl.eventHandlers[evTpType],
		func(ev *eventpb.Event) error {
			evTp := ev.Type.(EvWrapper)
			// The safety of this cast is verified by the runtime checks below.
			// This could be done much nicer and without runtime-casts if the protoc-generated wrappers exported a
			// function with a known name like Elem(), which would return the internal object.
			evTpPtr := (*evContainer[Ev])(unsafe.Pointer(&evTp))
			return handler((*evTpPtr).ev)
		})

	// These checks verify that an object of type EvWrapper can be safely interpreted as object of type evContainer[Ev].
	// They are only performed at the time of registration of the handler (which is supposedly at the very beginning
	// of the program execution) and make sure that no unexpected runtime errors will happen during later stages
	// of protocol execution and no memory will be corrupted by the pointer cast between (*EvWrapper) and
	// (*evContainer[Ev]). These checks may fail in case of a change to the protobuf API or when erroneous type
	// parameters (EvWrapper and Ev) are provided.
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

// UponCondition registers a special type of handler that will be invoked each time after processing a batch of events.
// The handler is assumed to represent a conditional action: it is supposed to check some predicate on the state
// and perform actions if the predicate evaluates is satisfied.
func UponCondition(m Module, handler func() error) {
	m.GetDslHandle().impl.conditionHandlers = append(m.GetDslHandle().impl.conditionHandlers, handler)
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (m *dslModuleImpl) ImplementsModule() {}

func EmitEvent(m Module, ev *eventpb.Event) {
	m.GetDslHandle().impl.outputEvents.PushBack(ev)
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
