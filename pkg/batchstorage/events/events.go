package events

import (
	bspb "github.com/filecoin-project/mir/pkg/pb/batchstoragepb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Event creates an eventpb.Event out of an availabilitypb.Event.
func Event(dest t.ModuleID, ev *bspb.Event) *eventpb.Event {
	return &eventpb.Event{
		DestModule: dest.Pb(),
		Type: &eventpb.Event_BatchStorage{
			BatchStorage: ev,
		},
	}
}

func StoreBatchOnQuorum(dest t.ModuleID, metadata []byte, origin *bspb.StoreBatchOnQuorumOrigin) *eventpb.Event {
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_Store{
			Store: &bspb.StoreBatchOnQuorum{
				Metadata: metadata,
				Origin:   origin,
			},
		},
	})
}

func BatchStoredOnQuorum(dest t.ModuleID, cert *bspb.RetrieveCert, origin *bspb.StoreBatchOnQuorumOrigin) *eventpb.Event {
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_Stored{
			Stored: &bspb.BatchStoredOnQuorum{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

func VerifyBatch(dest t.ModuleID, metadata []byte, requestID uint64) *eventpb.Event {
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_VerifyBatch{
			VerifyBatch: &bspb.VerifyBatch{
				Metadata:  metadata,
				RequestId: requestID,
			},
		},
	})
}

func BatchVerified(dest t.ModuleID, err error, requestID uint64) *eventpb.Event {
	valid, errStr := t.ErrorPb(err)
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_BatchVerified{
			BatchVerified: &bspb.BatchVerified{
				Valid:     valid,
				Err:       errStr,
				RequestId: requestID,
			},
		},
	})
}

func VerifyCert(dest t.ModuleID, cert *bspb.RetrieveCert, origin *bspb.VerifyCertOrigin) *eventpb.Event {
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_VerifyCert{
			VerifyCert: &bspb.VerifyCert{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

func CertVerified(dest t.ModuleID, err error, origin *bspb.VerifyCertOrigin) *eventpb.Event {
	valid, errStr := t.ErrorPb(err)
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_CertVerified{
			CertVerified: &bspb.CertVerified{
				Valid:  valid,
				Err:    errStr,
				Origin: origin,
			},
		},
	})
}

func RetrieveTransactions(dest t.ModuleID, cert *bspb.RetrieveCert, origin *bspb.RetrieveTransactionsOrigin) *eventpb.Event {
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_RetrieveTransactions{
			RetrieveTransactions: &bspb.RetrieveTransactions{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

func TransactionsRetrieved(dest t.ModuleID, txIDs []t.TxID, txs []*requestpb.Request, origin *bspb.RetrieveTransactionsOrigin) *eventpb.Event {
	return Event(dest, &bspb.Event{
		Type: &bspb.Event_TransactionsRetrieved{
			TransactionsRetrieved: &bspb.TransactionsRetrieved{
				TxIds:  t.TxIDSlicePb(txIDs),
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}
