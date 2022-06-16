package test

import (
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	"github.com/filecoin-project/mir/pkg/types"
	"testing"
)

func TestDslModule_ApplyEvents(test *testing.T) {
	m := dsl.NewModule()

	var sentecho = false
	var delivered = false
	var echos map[types.ModuleID][]byte

	dsl.UponEvent[eventpb.Event_Request](m, func(ev *requestpb.Request) (*events.EventList, error) {
		ev.GetData()
	})
}
