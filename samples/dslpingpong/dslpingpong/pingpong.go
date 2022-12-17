package dslpingpong

import (
	"fmt"
	"time"

	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	ppdsl "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb/dsl"
	ppevents "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb/events"
	ppmsgs "github.com/filecoin-project/mir/pkg/pb/dslpingpongpb/msgs"
	eventpbdsl "github.com/filecoin-project/mir/pkg/pb/eventpb/dsl"
	eventpbtypes "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

type Pingpong struct {
	ownID  t.NodeID
	nextSn uint64
}

// ModuleConfig sets the module ids. All participants are expected to use identical module configurations.
type ModuleConfig struct {
	Self      t.ModuleID // id of this module
	Timer     t.ModuleID // id of the timer module
	Transport t.ModuleID // id of the network module
}

// DefaultModuleConfig returns a valid module config with default names for all modules.
func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:      "pingpong",
		Timer:     "timer",
		Transport: "transport",
	}
}

func NewPingPong(mc *ModuleConfig, ownID t.NodeID) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	var nextSn t.SeqNr

	dsl.UponInit(m, func() error {
		nextSn = 0

		eventpbdsl.TimerRepeat(
			m,
			mc.Timer,
			[]*eventpbtypes.Event{ppevents.PingTime(mc.Self)},
			t.TimeDuration(time.Second),
			0,
		)

		return nil
	})

	ppdsl.UponPingTime(m, func() error {
		var destID t.NodeID
		if ownID == "0" {
			destID = "1"
		} else {
			destID = "0"
		}

		eventpbdsl.SendMessage(m, mc.Transport, ppmsgs.Ping(mc.Self, nextSn), []t.NodeID{destID})
		nextSn++

		return nil
	})

	ppdsl.UponPingReceived(m, func(from t.NodeID, seqNr t.SeqNr) error {
		fmt.Printf("Received ping from %s: %d\n", from, seqNr)
		eventpbdsl.SendMessage(m, mc.Transport, ppmsgs.Pong(mc.Self, seqNr), []t.NodeID{from})
		return nil
	})

	ppdsl.UponPongReceived(m, func(from t.NodeID, seqNr t.SeqNr) error {
		fmt.Printf("Received pong from %s: %d\n", from, seqNr)
		return nil
	})

	return m
}
