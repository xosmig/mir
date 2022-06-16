package cb

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	"github.com/filecoin-project/mir/pkg/pb/cbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/messagepb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

type ModuleConfig struct {
	Self   t.ModuleID
	Net    t.ModuleID
	Crypto t.ModuleID
}

func DefaultModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Self:   "cb",
		Net:    "net",
		Crypto: "crypto",
	}
}

type ModuleParams struct {
	AllNodes []t.NodeID
	Leader   t.NodeID
}

func (params *ModuleParams) GetN() int {
	return len(params.AllNodes)
}

func (params *ModuleParams) GetF() int {
	return (params.GetN() - 1) / 3
}

type cbModuleState struct {
	sentEcho  bool
	sentFinal bool
	delivered bool
	echoSigs  map[t.ModuleID][]byte
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeId t.NodeID) modules.PassiveModule {
	m := dsl.NewModule()

	state := cbModuleState{
		sentEcho:  false,
		sentFinal: false,
		delivered: false,
		echoSigs:  make(map[t.ModuleID][]byte),
	}

	// upon event <bcb, Broadcast | m> do
	dsl.UponEvent[eventpb.Event_Request](m, func(ev *requestpb.Request) error {
		if nodeId != params.Leader {
			return fmt.Errorf("only the leader node can receive client requests")
		}
		dsl.SendMessage(m, mc.Net, StartMessage(mc.Self, ev.Data), params.AllNodes)
		return nil
	})

	// upon event <al, Deliver | p, [Send, m]> such that p = s and sentecho = false do
	UponMessageReceived[cbpb.CBMessage_EchoMessage](m, func(msg *cbpb.EchoMessage) error {
		dsl.SignRequest(m, mc.Crypto, [][]byte{msg.Data}, SignRequestOrigin())
		return nil
	})

	dsl.UponEvent[eventpb.Event_SignResult](m, func(ev *eventpb.SignResult) error {
		dsl.SendMessage(m, mc.Net, EchoMessage(mc.Self, ev.Signature), []t.NodeID{params.Leader})
		return nil
	})

	dsl.UponEventWithCondition[eventpb.Event_MessageReceived](m, isEchoMessage, func(ev *eventpb.MessageReceived) error {
		msg := ev.Msg.GetCb().GetEchoMessage()
		sig := msg.GetSig()
		if len(state.echoSigs) >= params.GetN() - params.GetF() {
			dsl.VerifyNodeSigs(m, mc.Crypto, )
		}
	})



	dsl.UponOneShotCondition(m, func() bool { state. }, func() error {
		val := state.decided;
		Deliver(val)
	})
}

func isStartMessage(ev *eventpb.MessageReceived) bool {
	cbMessage, ok := ev.Msg.Type.(*messagepb.Message_Cb)
	if !ok {
		return false
	}

	_, ok = cbMessage.Cb.Type.(*cbpb.CBMessage_StartMessage)
	return ok
}

func isEchoMessage(ev *eventpb.MessageReceived) bool {
	cbMessage, ok := ev.Msg.Type.(*messagepb.Message_Cb)
	if !ok {
		return false
	}

	_, ok = cbMessage.Cb.Type.(*cbpb.CBMessage_EchoMessage)
	return ok
}

