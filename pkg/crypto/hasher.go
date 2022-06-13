package crypto

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/statuspb"
	t "github.com/filecoin-project/mir/pkg/types"
	"hash"
)

type HashImpl interface {
	New() hash.Hash
}

type Hasher struct {
	modules.Module

	hashImpl HashImpl
}

func NewHasher(hashImpl HashImpl) *Hasher {
	return &Hasher{
		hashImpl: hashImpl,
	}
}

func (hasher *Hasher) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {

	// TODO: Parallelize this.

	eventsOut := &events.EventList{}

	iter := eventsIn.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {
		evts, err := hasher.ApplyEvent(event)
		if err != nil {
			return nil, err
		}
		eventsOut.PushBackList(evts)
	}

	return eventsOut, nil
}

func (hasher *Hasher) ApplyEvent(event *eventpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *eventpb.Event_HashRequest:
		// HashRequest is the only event understood by the hasher module.

		// Create a slice for the resulting digests containing one element for each data item to be hashed.
		digests := make([][]byte, len(e.HashRequest.Data))

		// Hash each data item contained in the event
		for i, data := range e.HashRequest.Data {

			// One data item consists of potentially multiple byte slices.
			// Add each of them to the hash function.
			h := hasher.hashImpl.New()
			for _, d := range data.Data {
				h.Write(d)
			}

			// Save resulting digest in the result slice
			digests[i] = h.Sum(nil)
		}

		// Return all computed digests in one common event.
		return (&events.EventList{}).PushBack(
			events.HashResult(t.ModuleID(e.HashRequest.Origin.Module), digests, e.HashRequest.Origin),
		), nil
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected type of Hash event: %T", event.Type)
	}
}

func (hasher *Hasher) Status() (s *statuspb.ProtocolStatus, err error) {
	//TODO implement me
	panic("implement me")
}
