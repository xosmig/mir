package events

import (
	qspb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/quorumstoragepb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"github.com/filecoin-project/mir/pkg/pb/requestpb"
	t "github.com/filecoin-project/mir/pkg/types"
)

// Event creates an eventpb.Event out of an qspb.Event.
func Event(dest t.ModuleID, ev *qspb.Event) *eventpb.Event {
	return &eventpb.Event{
		DestModule: dest.Pb(),
		Type: &eventpb.Event_QuorumStorage{
			QuorumStorage: ev,
		},
	}
}

// StoreBatchOnQuorum can be used to create a batch and store in on a quorum of nodes wit the given metadata.
func StoreBatchOnQuorum(dest t.ModuleID, metadata []byte, origin *qspb.StoreBatchOnQuorumOrigin) *eventpb.Event {
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_Store{
			Store: &qspb.StoreBatchOnQuorum{
				Metadata: metadata,
				Origin:   origin,
			},
		},
	})
}

// BatchStoredOnQuorum is a response to a StoreBatchOnQuorum event.
func BatchStoredOnQuorum(dest t.ModuleID, cert *qspb.RetrieveCert, origin *qspb.StoreBatchOnQuorumOrigin) *eventpb.Event {
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_Stored{
			Stored: &qspb.BatchStoredOnQuorum{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

// VerifyBatch is used by the quorum storage to check the validity of the batch metadata.
func VerifyBatch(dest t.ModuleID, metadata []byte, origin *qspb.VerifyBatchOrigin) *eventpb.Event {
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_VerifyBatch{
			VerifyBatch: &qspb.VerifyBatch{
				Metadata: metadata,
				Origin:   origin,
			},
		},
	})
}

// BatchVerified is a response to a VerifyBatch event.
func BatchVerified(dest t.ModuleID, err error, origin *qspb.VerifyBatchOrigin) *eventpb.Event {
	valid, errStr := t.ErrorPb(err)
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_BatchVerified{
			BatchVerified: &qspb.BatchVerified{
				Valid:  valid,
				Err:    errStr,
				Origin: origin,
			},
		},
	})
}

// VerifyCert can be used to verify the validity of the retrieval certifcate.
func VerifyCert(dest t.ModuleID, cert *qspb.RetrieveCert, origin *qspb.VerifyCertOrigin) *eventpb.Event {
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_VerifyCert{
			VerifyCert: &qspb.VerifyCert{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

// CertVerified is a response to a VerifyCert event.
func CertVerified(dest t.ModuleID, err error, origin *qspb.VerifyCertOrigin) *eventpb.Event {
	valid, errStr := t.ErrorPb(err)
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_CertVerified{
			CertVerified: &qspb.CertVerified{
				Valid:  valid,
				Err:    errStr,
				Origin: origin,
			},
		},
	})
}

// RetrieveTransactions can be used to retrieve the transactions in a batch.
func RetrieveTransactions(dest t.ModuleID, cert *qspb.RetrieveCert, origin *qspb.RetrieveTransactionsOrigin) *eventpb.Event {
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_RetrieveTransactions{
			RetrieveTransactions: &qspb.RetrieveTransactions{
				Cert:   cert,
				Origin: origin,
			},
		},
	})
}

// TransactionsRetrieved is a response to a RetrieveTransactions event.
func TransactionsRetrieved(dest t.ModuleID, txIDs []t.TxID, txs []*requestpb.Request, origin *qspb.RetrieveTransactionsOrigin) *eventpb.Event {
	return Event(dest, &qspb.Event{
		Type: &qspb.Event_TransactionsRetrieved{
			TransactionsRetrieved: &qspb.TransactionsRetrieved{
				TxIds:  t.TxIDSlicePb(txIDs),
				Txs:    txs,
				Origin: origin,
			},
		},
	})
}
