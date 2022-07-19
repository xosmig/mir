package multisigcollector

import (
	"github.com/filecoin-project/mir/pkg/deploytest"
	"github.com/filecoin-project/mir/pkg/events"
	mpevents "github.com/filecoin-project/mir/pkg/mempool/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/modules/mockmodules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	mppb "github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	"github.com/filecoin-project/mir/pkg/types"
	"github.com/filecoin-project/mir/pkg/util/eventutil"
	"github.com/filecoin-project/mir/pkg/util/sliceutil"
	"github.com/golang/mock/gomock"
	"github.com/xosmig/placeholders"
	"testing"
)

func TestMultisigCollector_WithMockedMempool(t *testing.T) {
	testCases := map[string]struct {
		nNodes             int
		idxOfCrashedNodes  []int
		idxOfProposerNodes []int
	}{
		"1 request, 4 nodes, all correct": {
			nNodes:             4,
			idxOfCrashedNodes:  nil,
			idxOfProposerNodes: []int{0},
		},
		"1 request, 4 nodes, 1 crashed": {
			nNodes:             4,
			idxOfCrashedNodes:  []int{2},
			idxOfProposerNodes: []int{1},
		},
		"1 request, 7 nodes, all correct": {
			nNodes:             7,
			idxOfCrashedNodes:  nil,
			idxOfProposerNodes: []int{2},
		},
		"1 request, 7 nodes, 2 crashed": {
			nNodes:             7,
			idxOfCrashedNodes:  []int{0, 6},
			idxOfProposerNodes: []int{2},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			nodeIDs := deploytest.NewNodeIDs(tc.nNodes)

			mc := DefaultModuleConfig()

			params := &ModuleParams{
				InstanceUID: []byte("test instance"),
				AllNodes: nodeIDs,
			}

			var nodeModules map[types.NodeID]modules.Modules
			transportLayer := deploytest.NewFakeTransport(nodeIDs)

			crashedNodes := sliceutil.ToSet(tc.idxOfCrashedNodes)
			proposerNodes := sliceutil.ToSet(tc.idxOfProposerNodes)

			for i := 0; i < tc.nNodes; i++ {
				nodeID := nodeIDs[i]
				self := NewModule(mc, params, nodeID)

				mempool := mockmodules.NewMockPassiveModule(ctrl)
				if proposerNodes[i] {
					expectProposeBatch(mempool, mc)
				}

				var net modules.Module
				if !crashedNodes[i] {
					net = transportLayer.Link(nodeID)
				}

				nodeModules[nodeID] = modules.Modules{
					mc.Self: self,
					mc.Mempool: mempool,
					mc.Net:     net,
					mc.Crypto:  ,
				}
			}
		})
	}
}

func expectProposeBatch(mempool *mockmodules.MockPassiveModule, mc *ModuleConfig) {
	txIDs := []types.TxID{
		"tx one",
		types.TxID([]byte{213, 137, 116, 20, 237, 78, 55, 52, 124, 200}), // this is a random sequence of bytes
	}

	txs := [][]byte{
		[]byte("transaction one payload"),
		{33, 175, 88, 120, 223, 24, 113, 35, 227, 159}, // this is a random sequence of bytes
	}

	mempool.EXPECT().Event(mpevents.RequestBatch(mc.Mempool, placeholders.Make[*mppb.RequestBatchOrigin](t))).
		DoAndReturn(func(ev *eventpb.Event) (*events.EventList, error) {
			response := mpevents.NewBatch(mc.Self, txIDs, txs, eventutil.GetOrigin[mppb.RequestBatchOrigin](ev))
			return events.ListOf(response), nil
		})

	mempool.EXPECT().Event(mpevents.RequestBatchID(mc.Mempool, txIDs, placeholders.Make[*mppb.RequestBatchIDOrigin](t))).
		DoAndReturn(func(ev *eventpb.Event) (*events.EventList, error) {
			response := mpevents.BatchIDResponse(mc.Self, "test batch ID", eventutil.GetOrigin[mppb.RequestBatchIDOrigin](ev))
			return events.ListOf(response), nil
		})
}
