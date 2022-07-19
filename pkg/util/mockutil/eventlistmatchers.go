package mockutil

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/golang/mock/gomock"
	"reflect"
)

// EventListOf matches a list of events.
// Each element is either an object of type *eventpb.Event or a gomock matcher.
func EventListOf(evsOrMatchers ...any) gomock.Matcher {
	for _, evOrMatcher := range evsOrMatchers {
		_, isMatcher := evOrMatcher.(gomock.Matcher)
		_, isEvent := evOrMatcher.(*eventpb.Event)
		if !isEvent && !isMatcher {
			panic(fmt.Errorf("expected an event or a gomock.Matcher, got: %v", evOrMatcher))
		}
	}

	return &eventListOfMatcher{evsOrMatchers}
}

type eventListOfMatcher struct {
	evsOrMatchers []any
}

// Matches returns whether x is a match.
func (m *eventListOfMatcher) Matches(x any) bool {
	evList, ok := x.(*events.EventList)
	if !ok {
		return false
	}

	evs := evList.Slice()

	if len(evs) != len(m.evsOrMatchers) {
		return false
	}

	for i := range evs {
		if !MatchEvent(m.evsOrMatchers[i], evs[i]) {
			return false
		}
	}

	return true
}

// String describes what the matcher matches.
func (m *eventListOfMatcher) String() string {
	return fmt.Sprintf("{event list %v}", m.evsOrMatchers)
}

// MatchEvent checks if the given event matches the expectation.
func MatchEvent(expectedEvOrMatcher any, ev *eventpb.Event) bool {
	if matcher, ok := expectedEvOrMatcher.(gomock.Matcher); ok {
		return matcher.Matches(ev)
	}

	if expectedEv, ok := expectedEvOrMatcher.(*eventpb.Event); ok {
		return reflect.DeepEqual(ev, expectedEv)
	}

	panic(fmt.Errorf("expected an event or a gomock.Matcher, got: %v", ev))
}
