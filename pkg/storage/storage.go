package storage

import (
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"sync"
)

type Storage interface {
	// GetState reads a value for a specific key. The value may be read from cache
	// and not persistent yet.
	GetState(key string) ([]byte, error)

	// Write writes the transaction batch and updated state to underlying storage.
	Write(batchCommittedResult *BatchCommittedResult)
}

type BatchCommittedResult struct {
	TransactionContexts []*transaction.Context
	BatchNum            int64
	WrittenSignal       chan struct{}
	SyncedSignal        chan struct{}
}

type MemStorage struct {
	store sync.Map
}

func (ms *MemStorage) GetState(key string) ([]byte, error) {
	value, _ := ms.store.Load(key)
	if value == nil {
		return nil, nil
	}
	return value.([]byte), nil
}

func (ms *MemStorage) Write(batchCommittedResult *BatchCommittedResult) {
	for _, tctx := range batchCommittedResult.TransactionContexts {
		for _, kvWrite := range tctx.RWSet.WSet {
			if kvWrite.IsDelete {
				ms.store.Delete(kvWrite.Key)
			} else {
				ms.store.Store(kvWrite.Key, kvWrite.Value)
			}
		}
	}
	if batchCommittedResult.WrittenSignal != nil {
		batchCommittedResult.WrittenSignal <- struct{}{}
	}

	if batchCommittedResult.SyncedSignal != nil {
		batchCommittedResult.SyncedSignal <- struct{}{}
	}
}
