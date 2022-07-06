package submodules

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
)

// SubmoduleID is used to identify submodules of a module
type SubmoduleID string

// Module is an extension of dsl.Module that supports submodules functionality.
type Module interface {
	dsl.Module
	SubmodulesHandle() Handle
}

// Handle is used internally to manage the submodules.
type Handle struct {
	impl *moduleImpl
}

type moduleImpl struct {
	dsl.Module
	submodules    map[SubmoduleID]modules.PassiveModule
	pendingEvents map[SubmoduleID][]*eventpb.Event
}

// SubmodulesHandle is used internally to manage the submodules.
func (m *moduleImpl) SubmodulesHandle() Handle {
	return Handle{m}
}

func AddSubmodules(underlying dsl.Module, onUnknownSubmodule func(id SubmoduleID, ev *eventpb.Event) error) Module {
	if len(dsl.GetEventHandlers[*eventpb.Event_Submodule](underlying.DslHandle())) > 0 {
		panic(fmt.Errorf("module '%v' already has submodules enabled", underlying.ModuleID()))
	}

	m := &moduleImpl{
		Module:     underlying,
		submodules: make(map[SubmoduleID]modules.PassiveModule),
	}

	dsl.UponEvent[*eventpb.Event_Submodule](m, func(ev *eventpb.SubmoduleEvent) error {
		submoduleID := SubmoduleID(ev.SubmoduleId)
		_, ok := m.submodules[submoduleID]
		if !ok {
			err := onUnknownSubmodule(submoduleID, ev.Event)
			if err != nil {
				return err
			}
			// check if the handler created the submodule
			_, ok = m.submodules[submoduleID]
			if !ok {
				return nil
			}
		}

		m.pendingEvents[submoduleID] = append(m.pendingEvents[submoduleID], ev.Event)
		return nil
	})

	dsl.UponCondition(m, func() error {
		if len(m.pendingEvents) > 0 {
			for submoduleID, evs := range m.pendingEvents {
				eventsOut, err := m.submodules[submoduleID].ApplyEvents(events.ListOf(evs...))
				if err != nil {
					return err
				}

				iter := eventsOut.Iterator()
				for ev := iter.Next(); ev != nil; ev = iter.Next() {
					dsl.EmitEvent(m, ev)
				}
			}
		}
		return nil
	})

	return m
}

func Start(m Module, id SubmoduleID, submodule modules.PassiveModule) {
	if _, ok := m.SubmodulesHandle().impl.submodules[id]; ok {
		panic(fmt.Errorf("module %v alread has a submodule with id %v", m.ModuleID(), id))
	}
	m.SubmodulesHandle().impl.submodules[id] = submodule
}

func Stop(m Module, id SubmoduleID) {
	if _, ok := m.SubmodulesHandle().impl.submodules[id]; !ok {
		panic(fmt.Errorf("module %v does not have a submodule with id %v", m.ModuleID(), id))
	}
	delete(m.SubmodulesHandle().impl.submodules, id)
}

func UponEventForUnknownSubmodule(m Module, id SubmoduleID, ev *eventpb.Event) error {

}
