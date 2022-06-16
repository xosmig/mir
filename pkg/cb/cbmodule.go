package cb

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/cb/cbdsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	"github.com/filecoin-project/mir/pkg/pb/cbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
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
	echos     map[t.NodeID][]byte
	echoSigs  map[t.NodeID][]byte
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeId t.NodeID) modules.PassiveModule {
	m := cbdsl.NewModule(mc.Self)

	state := cbModuleState{
		sentEcho:  false,
		sentFinal: false,
		delivered: false,
		echoSigs:  make(map[t.NodeID][]byte),
	}

	// upon event <bcb, Broadcast | m> do    // only process s
	dsl.UponRequest(m, func(clientId string, reqNo uint64, data []byte, authenticator []byte) error {
		if nodeId != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		dsl.SendMessage(m, mc.Net, StartMessage(mc.Self, data), params.AllNodes)
		return nil
	})

	// upon event <al, Deliver | p, [Send, m]> such that p = s and sentecho = false do
	cbdsl.UponStartMessageReceived(m, func(from t.NodeID, msg *cbpb.StartMessage) error {
		if from == params.Leader && state.sentEcho {
			state.sentEcho = true
			cbdsl.SignRequest(m, mc.Crypto, [][]byte{msg.Data})
		}
		return nil
	})

	dsl.UponSignResult(m, func(signature []byte, _ *eventpb.SignOrigin) error {
		dsl.SendMessage(m, mc.Net, EchoMessage(mc.Self, signature), []t.NodeID{params.Leader})
		return nil
	})

	// upon event <al, Deliver | p, [ECHO, m, σ]> do    // only process s
	cbdsl.UponEchoMessageReceived(m, func(from t.NodeID, msg *cbpb.EchoMessage) error {
		// if echos[p] = ⊥ ∧ verifysig(p, bcb||p||ECHO||m, σ) then
		_, alreadyReceived := state.echos[from]
		if nodeId == params.Leader && !alreadyReceived {
			state.echos[from] = msg.Data
			cbdsl.VerifyNodeSignature(m, mc.Crypto, [][]byte{msg.Data}, msg.Signature, from)
		}
		return nil
	})

	dsl.
		dsl.UponOneShotCondition(m,
		func() bool { return len(state.echos) >= params.GetN()-params.GetF() && !state.sentFinal },
		func() error {
			val := state.decided
			Deliver(val)
		})
}

//func isStartMessage(ev *eventpb.MessageReceived) bool {
//	cbMessage, ok := ev.Msg.Type.(*messagepb.Message_Cb)
//	if !ok {
//		return false
//	}
//
//	_, ok = cbMessage.Cb.Type.(*cbpb.CBMessage_StartMessage)
//	return ok
//}
//
//func isEchoMessage(ev *eventpb.MessageReceived) bool {
//	cbMessage, ok := ev.Msg.Type.(*messagepb.Message_Cb)
//	if !ok {
//		return false
//	}
//
//	_, ok = cbMessage.Cb.Type.(*cbpb.CBMessage_EchoMessage)
//	return ok
//}
//
