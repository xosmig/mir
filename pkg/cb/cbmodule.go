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

	sentEcho         bool
	sentFinal        bool
	delivered        bool
	receivedEcho     map[t.NodeID]bool
	echoSigs         map[t.NodeID][]byte
	verifiedEchoSigs map[t.NodeID][]byte
}

// NewModule returns a passive module for the Signed Echo Broadcast from the textbook "Introduction to reliable and
// secure distributed programming". It serves as a motivating example for the DSL module interface.
func NewModule(mc *ModuleConfig, params *ModuleParams, nodeId t.NodeID) modules.PassiveModule {
	m := cbdsl.NewModule(mc.Self)

	state := cbModuleState{
		request:          nil,
		sentEcho:         false,
		sentFinal:        false,
		delivered:        false,
		receivedEcho:     make(map[t.NodeID]bool),
		echoSigs:         make(map[t.NodeID][]byte),
		verifiedEchoSigs: make(map[t.NodeID][]byte),
	}

	// upon event <bcb, Broadcast | m> do    // only process s
	dsl.UponRequest(m, func(clientId t.ClientID, reqNo uint64, data []byte, authenticator []byte) error {
		if nodeId != params.Leader {
			return fmt.Errorf("only the leader node can receive requests")
		}
		state.request = data
		dsl.SendMessage(m, mc.Net, StartMessage(mc.Self, data), params.AllNodes)
		return nil
	})

	// upon event <al, Deliver | p, [Send, m]> ...
	cbdsl.UponStartMessageReceived(m, func(from t.NodeID, msg *cbpb.StartMessage) error {
		// ... such that p = s and sentecho = false do
		if from == params.Leader && !state.sentEcho {
			cbdsl.SignRequest(m, mc.Crypto, [][]byte{msg.Data})
		}
		return nil
	})

	dsl.UponSignResult(m, func(signature []byte, origin *eventpb.SignOrigin) error {
		if !state.sentEcho {
			state.sentEcho = true
			dsl.SendMessage(m, mc.Net, EchoMessage(mc.Self, signature), []t.NodeID{params.Leader})
		}
		return nil
	})

	// upon event <al, Deliver | p, [ECHO, m, σ]> do    // only process s
	cbdsl.UponEchoMessageReceived(m, func(from t.NodeID, msg *cbpb.EchoMessage) error {
		// if echos[p] = ⊥ ∧ verifysig(p, bcb||p||ECHO||m, σ) then
		if nodeId == params.Leader && !state.receivedEcho[from] && state.request != nil {
			state.receivedEcho[from] = true
			state.echoSigs[from] = msg.Signature
			cbdsl.VerifyNodeSignature(m, mc.Crypto, [][]byte{state.request}, msg.Signature, from)
		}
		return nil
	})

	cbdsl.UponNodeSignatureVerified(m, func(nodeId t.NodeID, valid bool, err error) error {
		if valid && err == nil {
			state.verifiedEchoSigs[nodeId] = state.echoSigs[nodeId]
		}
		return nil
	})

	// upon exists m != ⊥ such that #({p ∈ Π | echos[p] = m}) > (N+f)/2 and sentfinal = FALSE do
	dsl.UponCondition(m, func() error {
		if len(state.verifiedEchoSigs) > (params.GetN()+params.GetF())/2 && !state.sentFinal {
			state.sentFinal = true
			certSigners := maputil.GetKeys(state.verifiedEchoSigs)
			certSignatures := maputil.GetValues(state.verifiedEchoSigs)
			dsl.SendMessage(m, mc.Net,
				FinalMessage(mc.Self, state.request, certSigners, certSignatures),
				params.AllNodes)
		}
		return nil
	})

	// upon event <al, Deliver | p, [FINAL, m, Σ]> do
	cbdsl.UponFinalMessageReceived(m, func(from t.NodeID, msg *cbpb.FinalMessage) error {
		// if #({p ∈ Π | Σ[p] != ⊥ ∧ verifysig(p, bcb||p||ECHO||m, Σ[p])}) > (N+f)/2 and delivered = FALSE do
		if len(msg.Signers) == len(msg.Signatures) && len(msg.Signers) > (params.GetN()+params.GetF())/2 && !state.delivered {

		}
		return nil
	})

	return m.GetPassiveModule()
}
