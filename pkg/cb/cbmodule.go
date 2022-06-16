package cb

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/cb/cbdsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/modules/dsl"
	"github.com/filecoin-project/mir/pkg/pb/cbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/maputil"
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
	// this variable is not part of the original protocol description, but it greatly simplifies the code
	request []byte

	sentEcho     bool
	sentFinal    bool
	delivered    bool
	receivedEcho map[t.NodeID]bool
	echoSigs     map[t.NodeID][]byte
}

func NewModule(mc *ModuleConfig, params *ModuleParams, nodeId t.NodeID) modules.PassiveModule {
	m := cbdsl.NewModule(mc.Self)

	state := cbModuleState{
		request:      nil,
		sentEcho:     false,
		sentFinal:    false,
		delivered:    false,
		receivedEcho: make(map[t.NodeID]bool),
		echoSigs:     make(map[t.NodeID][]byte),
	}

	// upon event <bcb, Broadcast | m> do    // only process s
	dsl.UponRequest(m, func(clientId string, reqNo uint64, data []byte, authenticator []byte) error {
		if nodeId != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		state.request = data
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
		if nodeId == params.Leader && !state.receivedEcho[from] && state.request != nil {
			state.receivedEcho[from] = true
			dsl.VerifyNodeSignature(m, mc.Crypto, [][]byte{state.request}, msg.Signature, from)
		}
		return nil
	})

	dsl.UponNodeSignatureVerified(m, func(TODO) error {
		TODO
	})

	// upon exists m != ⊥ such that #({p ∈ Π | echos[p] = m}) > (N + f) / 2 and sentfinal = FALSE do
	dsl.UponCondition(m, func() error {
		if len(state.echoSigs) > (params.GetN()+params.GetF())/2 && !state.sentFinal {
			state.sentFinal = true
			dsl.SendMessage(m, mc.Net,
				FinalMessage(mc.Self, state.request, maputil.GetKeys(state.echoSigs), maputil.GetValues(state.echoSigs)),
				params.AllNodes)
		}
		return nil
	})

	// upon event <al, Deliver | p, [FINAL, m, Σ]> do
	cbdsl.UponFinalMessageReceived(m, func(from t.NodeID, msg *cbpb.FinalMessage) error {
		TODO
	})

	return m.GetPassiveModule()
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
