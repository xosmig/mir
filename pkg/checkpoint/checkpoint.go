package checkpoint

import (
	"fmt"

	"github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/pb/checkpointpb"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	"github.com/filecoin-project/mir/pkg/serializing"
	t "github.com/filecoin-project/mir/pkg/types"
)

// StableCheckpoint represents a stable checkpoint.
type StableCheckpoint checkpointpb.StableCheckpoint

// SeqNr returns the sequence number of the stable checkpoint.
// It is defined as the number of sequence numbers comprised in the checkpoint, or, in other words,
// the first (i.e., lowest) sequence number not included in the checkpoint.
func (sc *StableCheckpoint) SeqNr() t.SeqNr {
	return t.SeqNr(sc.Sn)
}

// Memberships returns the memberships configured for the epoch of this checkpoint
// and potentially several subsequent ones.
func (sc *StableCheckpoint) Memberships() []map[t.NodeID]t.NodeAddress {
	return t.MembershipSlice(sc.Snapshot.Configuration.Memberships)
}

// Epoch returns the epoch associated with this checkpoint.
// It is the epoch **started** by this checkpoint, **not** the last one included in it.
func (sc *StableCheckpoint) Epoch() t.EpochNr {
	return t.EpochNr(sc.Snapshot.Configuration.EpochNr)
}

// StateSnapshot returns the serialized application state and system configuration associated with this checkpoint.
func (sc *StableCheckpoint) StateSnapshot() *commonpb.StateSnapshot {
	return sc.Snapshot
}

// Pb returns a protobuf representation of the stable checkpoint.
func (sc *StableCheckpoint) Pb() *checkpointpb.StableCheckpoint {
	return (*checkpointpb.StableCheckpoint)(sc)
}

// VerifyCert verifies the certificate of the stable checkpoint using the provided hash implementation and verifier.
// The same (or corresponding) modules must have been used when the certificate was created by the checkpoint module.
// The has implementation is a crypto.HashImpl used to create a Mir hasher module and the verifier interface
// is a subset of the crypto.Crypto interface (narrowed down to only the Verify function).
// Thus, the same (or equivalent) crypto implementation that was used to create checkpoint
// can be used as a Verifier to verify it.
//
// Note that VerifyCert performs all the necessary hashing and signature verifications synchronously
// (only returns when the signature is verified). This may become a very computationally expensive operation.
// It is thus recommended not to use this function directly within a sequential protocol implementation,
// and rather delegating the hashing and signature verification tasks
// to dedicated modules using the corresponding events.
// Also, in case the verifier implementation is used by other goroutines,
// make sure that calling Vetify on it is thread-safe.
//
// For simplicity, we require all nodes that signed the certificate to be contained in the provided membership,
// as well as all signatures to be valid.
// Moreover, the number of nodes that signed the certificate must be greater than one third of the membership size.
func (sc *StableCheckpoint) VerifyCert(h crypto.HashImpl, v Verifier, membership map[t.NodeID]t.NodeAddress) error {

	// Check if there is enough signatures.
	n := len(membership)
	f := (n - 1) / 3
	if len(sc.Cert) <= f+1 {
		return fmt.Errorf("not enough signatures in certificate: got %d, expected more than %d",
			len(sc.Cert), f+1)
	}

	// Check whether all signatures are valid.
	snapshotData := serializing.SnapshotForHash(sc.StateSnapshot())
	snapshotHash := hash(snapshotData, h)
	sigData := serializing.CheckpointForSig(sc.Epoch(), sc.SeqNr(), snapshotHash)
	for nodeID, sig := range sc.Cert {
		// For each signature in the certificate...

		// Check if the signing node is also in the given membership, thus "authorized" to sign.
		// TODO: Once nodes are identified by more than their ID
		//   (e.g., if a separate putlic key is part of their identity), adapt the check accordingly.
		if _, ok := membership[t.NodeID(nodeID)]; !ok {
			return fmt.Errorf("node %v not in membership", nodeID)
		}

		// Check if the signature is valid.
		if err := v.Verify(sigData, sig, t.NodeID(nodeID)); err != nil {
			return fmt.Errorf("signature verification error (node %v): %w", nodeID, err)
		}
	}
	return nil
}

func Genesis(initialStateSnapshot *commonpb.StateSnapshot) *StableCheckpoint {
	return &StableCheckpoint{
		Sn:       0,
		Snapshot: initialStateSnapshot,
		Cert:     map[string][]byte{},
	}
}

// The Verifier interface represents a subset of the crypto.Crypto interface
// that can be used for verifying stable checkpoint certificates.
type Verifier interface {
	// Verify verifies a signature produced by the node with ID nodeID over data.
	// Returns nil on success (i.e., if the given signature is valid) and a non-nil error otherwise.
	Verify(data [][]byte, signature []byte, nodeID t.NodeID) error
}

func hash(data [][]byte, hasher crypto.HashImpl) []byte {
	h := hasher.New()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
}
