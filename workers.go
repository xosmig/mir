/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0

Refactored: 1
*/

package mir

import (
	"context"
	"fmt"
	"github.com/filecoin-project/mir/pkg/modules"
	"runtime/debug"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/statuspb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Input and output channels for the modules within the Node.
// the Node.process() method reads and writes events
// to and from these channels to rout them between the Node's modules.
type workChans struct {

	// There is one channel per module to feed events into the module.
	protocol chan *events.EventList
	net      chan *events.EventList
	timer    chan *events.EventList

	// All modules write their output events in a common channel, from where the node processor reads and redistributes
	// the events to their respective workItems buffers.
	// External events are also funneled through this channel towards the workItems buffers.
	workItemInput chan *events.EventList

	// Events received during debugging through the Node.Step function are written to this channel
	// and inserted in the event loop.
	debugIn chan *events.EventList

	// During debugging, Events that would normally be inserted in the workItems event buffer
	// (and thus inserted in the event loop) are written to this channel instead if it is not nil.
	// If this channel is nil, those Events are discarded.
	debugOut chan *events.EventList

	genericWorkChans map[t.ModuleID]chan *events.EventList
}

// Allocate and return a new workChans structure.
func newWorkChans(modules *modules.Modules) workChans {
	genericWorkChans := make(map[t.ModuleID]chan *events.EventList)

	for moduleID := range modules.GenericModules {
		genericWorkChans[moduleID] = make(chan *events.EventList)
	}

	return workChans{
		protocol: make(chan *events.EventList),
		net:      make(chan *events.EventList),
		timer:    make(chan *events.EventList),

		workItemInput: make(chan *events.EventList),

		debugIn:  make(chan *events.EventList),
		debugOut: make(chan *events.EventList),

		genericWorkChans: genericWorkChans,
	}
}

// A function type used for performing the work of a module.
// It usually reads events from a work channel and writes the output to another work channel.
// Any error that occurs while performing the work is returned.
// When ctx is canceled, the function should return ErrStopped
type workFunc func(ctx context.Context) error

// Calls the passed work function repeatedly in an infinite loop until the work function returns an non-nil error.
// doUntilErr then sets the error in the Node's workErrNotifier and returns.
func (n *Node) doUntilErr(ctx context.Context, work workFunc) {
	for {
		err := work(ctx)
		if err != nil {
			n.workErrNotifier.Fail(err)
			return
		}
	}
}

// eventProcessor defines the type of the function that processes a single input events.EventList,
// producing a single output events.EventList.
// There is one such function defined for each Module that is executed in a loop by a worker goroutine.
type eventProcessor func(context.Context, *events.EventList) (*events.EventList, error)

// processEvents reads a single list of input Events from a work channel, strips off all associated follow-up Events,
// and processes the bare content of the list using the passed processing function.
// processEvents writes all the stripped off follow-up events along with any Events generated by the processing
// to the workItemInput channel, so they will be added to the workItems buffer for further processing.
//
// If the Node is configured to use an Interceptor, after having removed all follow-up Events,
// processEvents passes the list of input Events to the Interceptor.
//
// If exitC is closed, returns ErrStopped.
func (n *Node) processEvents(
	ctx context.Context,
	processFunc eventProcessor,
	eventSource <-chan *events.EventList,
) error {
	var eventsIn *events.EventList

	// Read input.
	select {
	case eventsIn = <-eventSource:
	case <-ctx.Done():
		return ErrStopped
	}

	// Remove follow-up Events from the input EventList,
	// in order to re-insert them in the processing loop after the input events have been processed.
	plainEvents, followUps := eventsIn.StripFollowUps()

	// Intercept the (stripped of all follow-ups) events that are about to be processed.
	// This is only for debugging / diagnostic purposes.
	n.interceptEvents(plainEvents)

	// Process events.
	newEvents, err := processFunc(ctx, plainEvents)
	if err != nil {
		return fmt.Errorf("could not process events: %w", err)
	}

	// Merge the pending follow-up Events with the newly generated Events.
	out := followUps.PushBackList(newEvents)

	// Return if no output was generated.
	// This is only an optimization to prevent the processor loop from handling empty EventLists.
	if out.Len() == 0 {
		return nil
	}

	// Write output.
	select {
	case n.workChans.workItemInput <- out:
	case <-ctx.Done():
		return ErrStopped
	}

	return nil
}

// processEventsPassive reads a single list of input Events from a work channel,
// strips off all associated follow-up Events,
// and processes the bare content of the list using the passed PassiveModule.
// processEventsPassive writes all the stripped off follow-up events along with any Events generated by the processing
// to the workItemInput channel, so they will be added to the workItems buffer for further processing.
//
// If the Node is configured to use an Interceptor, after having removed all follow-up Events,
// processEventsPassive passes the list of input Events to the Interceptor.
//
// If any error occurs, it is returned as the first parameter.
// If context is canceled, processEventsPassive might return a nil error with or without performing event processing.
// The second return value being true indicates that processing can continue
// and processEventsPassive should be called again.
// If the second return is false, processing should be terminated and processEventsPassive should not be called again.
func (n *Node) processEventsPassive(
	ctx context.Context,
	module modules.PassiveModule,
	eventSource <-chan *events.EventList,
) (error, bool) {
	var eventsIn *events.EventList
	var inputOpen bool

	// Read input.
	select {
	case eventsIn, inputOpen = <-eventSource:
		if !inputOpen {
			return nil, false
		}
	case <-ctx.Done():
		return nil, false
	case <-n.workErrNotifier.ExitC():
		return nil, false
	}

	// Remove follow-up Events from the input EventList,
	// in order to re-insert them in the processing loop after the input events have been processed.
	plainEvents, followUps := eventsIn.StripFollowUps()

	// Intercept the (stripped of all follow-ups) events that are about to be processed.
	// This is only for debugging / diagnostic purposes.
	n.interceptEvents(plainEvents)

	// Process events.
	newEvents, err := safelyApplyEvents(module, plainEvents)
	if err != nil {
		return err, false
	}

	// Merge the pending follow-up Events with the newly generated Events.
	out := followUps.PushBackList(newEvents)

	// Return if no output was generated.
	// This is only an optimization to prevent the processor loop from handling empty EventLists.
	if out.Len() == 0 {
		return nil, true
	}

	// Write output.
	select {
	case n.workChans.workItemInput <- out:
	case <-ctx.Done():
		return nil, false
	case <-n.workErrNotifier.ExitC():
		return nil, false
	}

	return nil, true
}

func safelyApplyEvents(
	module modules.PassiveModule,
	events *events.EventList,
) (result *events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = fmt.Errorf("module panicked: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = fmt.Errorf("module panicked: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return module.ApplyEvents(events)
}

// Module-specific wrappers for Node.ProcessEvents,
// associating each Module's processing function with its corresponding work channel.
// On top of that, the Protocol processing wrapper additionally sets the Node's exit status when done.

func (n *Node) doSendingWork(ctx context.Context) error {
	return n.processEvents(ctx, n.processSendEvents, n.workChans.net)
}

func (n *Node) doProtocolWork(ctx context.Context) (err error) {
	// On returning, sets the exit status of the protocol state machine in the work error notifier.
	defer func() {
		if err != nil {
			s, err := n.modules.Protocol.Status()
			n.workErrNotifier.SetExitStatus(&statuspb.NodeStatus{Protocol: s}, err)
			// TODO: Clean up status-related code.
		}
	}()
	return n.processEvents(ctx, n.processProtocolEvents, n.workChans.protocol)
}

func (n *Node) doTimerWork(ctx context.Context) (err error) {

	// Unlike other event processors that simply transform an event list to another event list,
	// the processor for the timer module needs direct access to the workItemsInput channel,
	// as the events it produces are (by definition) not available immediately.
	return n.processEvents(ctx, func(ctx context.Context, events *events.EventList) (*events.EventList, error) {
		return n.processTimerEvents(ctx, events, n.workChans.workItemInput)
	}, n.workChans.timer)
}

// TODO: Document the functions below.

func (n *Node) processSendEvents(_ context.Context, eventsIn *events.EventList) (*events.EventList, error) {
	eventsOut := &events.EventList{}

	iter := eventsIn.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {

		switch e := event.Type.(type) {
		case *eventpb.Event_SendMessage:
			for _, destID := range e.SendMessage.Destinations {
				if t.NodeID(destID) == n.ID {
					eventsOut.PushBack(events.MessageReceived(n.ID, e.SendMessage.Msg))
				} else {
					if err := n.modules.Net.Send(t.NodeID(destID), e.SendMessage.Msg); err != nil { // nolint
						// TODO: Handle sending errors (and remove "nolint" comment above).
					}
				}
			}
		default:
			return nil, fmt.Errorf("unexpected type of Net event: %T", event.Type)
		}
	}

	return eventsOut, nil
}

func (n *Node) processProtocolEvents(_ context.Context, eventsIn *events.EventList) (*events.EventList, error) {
	eventsOut := &events.EventList{}
	iter := eventsIn.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {

		newEvents, err := n.safeApplyProtocolEvent(event)
		if err != nil {
			return nil, fmt.Errorf("error applying protocol event: %w", err)
		}
		eventsOut.PushBackList(newEvents)
	}

	return eventsOut, nil
}

func (n *Node) safeApplyProtocolEvent(event *eventpb.Event) (result *events.EventList, err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = fmt.Errorf("panic in protocol state machine: %w\nStack trace:\n%s", rErr, string(debug.Stack()))
			} else {
				err = fmt.Errorf("panic in protocol state machine: %v\nStack trace:\n%s", r, string(debug.Stack()))
			}
		}
	}()

	return n.modules.Protocol.ApplyEvent(event), nil
}

// processTimerEvents processes the events destined to the timer module.
// Unlike other event processors, processTimerEvents does not return a list of resulting events
// based on the input event list, since those are (by definition) delayed.
// The returned EventList is thus always empty.
// Instead, processTimerEvents receives an additional channel (notifyChan)
// where its outputs are written at the appropriate times.
func (n *Node) processTimerEvents(
	ctx context.Context,
	eventsIn *events.EventList,
	notifyChan chan<- *events.EventList,
) (*events.EventList, error) {

	iter := eventsIn.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {

		// Based on event type, invoke the appropriate Timer function.
		// Note that events that later return to the event loop need to be copied in order to prevent a race condition
		// when they are later stripped off their follow-ups, as this happens potentially concurrently
		// with the original event being processed by the interceptor.
		switch e := event.Type.(type) {
		case *eventpb.Event_TimerDelay:
			n.modules.Timer.Delay(
				ctx,
				(&events.EventList{}).PushBackSlice(e.TimerDelay.Events),
				t.TimeDuration(e.TimerDelay.Delay),
				notifyChan,
			)
		case *eventpb.Event_TimerRepeat:
			n.modules.Timer.Repeat(
				ctx,
				(&events.EventList{}).PushBackSlice(e.TimerRepeat.Events),
				t.TimeDuration(e.TimerRepeat.Delay),
				t.TimerRetIndex(e.TimerRepeat.RetentionIndex),
				notifyChan,
			)
		case *eventpb.Event_TimerGarbageCollect:
			n.modules.Timer.GarbageCollect(t.TimerRetIndex(e.TimerGarbageCollect.RetentionIndex))
		default:
			return nil, fmt.Errorf("unexpected type of Timer event: %T", event.Type)
		}
	}

	return &events.EventList{}, nil
}
