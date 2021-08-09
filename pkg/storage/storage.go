package storage

import (
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type Storage interface {
	// GetUnstableState reads a value for a specific key. The value may be read from cache
	// and not persistent yet.
	GetUnstableState(key string) ([]byte, error)

	// GetStableState reads a value for a specific key. The value only read from non-volatile
	// storage.
	GetStableState(key string) ([]byte, error)

	// Write writes the transaction batch and updated state to underlying storage.
	Write(batchAndWSetSyncer *transaction.BatchAndWSetSyncer)
}
