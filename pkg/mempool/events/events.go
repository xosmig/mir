package events

import (
	"github.com/filecoin-project/mir/pkg/mempool"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	mppb "github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

func Event(dest t.ModuleID, ev *mppb.Event) *eventpb.Event {
	return &eventpb.Event{
		DestModule: dest.Pb(),
		Type: &eventpb.Event_Mempool{
			Mempool: ev,
		},
	}
}

func RequestBatch(dest t.ModuleID, origin *mppb.RequestBatchOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_RequestBatch{
			RequestBatch: &mppb.RequestBatch{
				Origin: origin,
			},
		},
	})
}

func NewBatch(dest t.ModuleID, txIDs []mempool.TxID, txs [][]byte, origin *mppb.RequestBatchOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_NewBatch{
			NewBatch: &mppb.NewBatch{
				TxIds:  txIDs.Pb(),
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}

func RequestTransactions(dest t.ModuleID, txIDs []string, origin *mppb.RequestTransactionsOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_RequestTransactions{
			RequestTransactions: &mppb.RequestTransactions{
				TxIds:  txIDs,
				Origin: origin,
			},
		},
	})
}

func TransactionsResponse(dest t.ModuleID, txs [][]byte, origin *mppb.RequestTransactionsOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_TransactionsResponse{
			TransactionsResponse: &mppb.TransactionsResponse{
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}

func RequestTransactionIDs(dest t.ModuleID, txs [][]byte, origin *mppb.RequestTransactionIDsOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_RequestTransactionIds{
			RequestTransactionIds: &mppb.RequestTransactionIDs{
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}

func TransactionsIDsResponse(dest t.ModuleID, txIDs []string, origin *mppb.RequestTransactionIDsOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_TransactionIdsResponse{
			TransactionIdsResponse: &mppb.TransactionIDsResponse{
				TxIds:  txIDs,
				Origin: origin,
			},
		},
	})
}

func RequestBatchID(dest t.ModuleID, txs [][]byte, origin *mppb.RequestBatchIDOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_RequestBatchId{
			RequestBatchId: &mppb.RequestBatchID{
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}

func BatchIDResponse(dest t.ModuleID, txIDs []string, origin *mppb.RequestBatchIDOrigin) *eventpb.Event {
	return Event(dest, &mppb.Event{
		Type: &mppb.Event_TransactionIdsResponse{
			TransactionIdsResponse: &mppb.TransactionIDsResponse{
				TxIds:  txIDs,
				Origin: origin,
			},
		},
	})
}
