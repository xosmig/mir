package common

import (
	msc "github.com/filecoin-project/mir/pkg/availability/multisigcollector"
)

// TxHash is a wrapper around the hash of the batch to be used as the map key.
// The underlying type is string and not []byte because []byte is not comparable and cannot be used as a map key.
type TxHash string

// BatchHash is a wrapper around the hash of the batch to be used as the map key.
// The underlying type is string and not []byte because []byte is not comparable and cannot be used as a map key.
type BatchHash string

type State struct {
	BatchStore       map[BatchHash][]TxHash
	TransactionStore map[TxHash][]byte
}

func SigMessage(instanceUID msc.InstanceUID, batchHash BatchHash) [][]byte {
	return [][]byte{instanceUID.Bytes(), []byte("BATCH_STORED"), []byte(batchHash)}
}
