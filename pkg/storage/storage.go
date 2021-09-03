package storage

import (
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"sync"
)

type Storage interface {
	// GetState reads a value for a specific key. The value may be read from cache
	// and not persistent yet.
	GetState(ns, key string) ([]byte, error)

	// Write writes the transaction batch and updated state to underlying storage.
	Write(batchCommittedResult *BatchCommittedResult)
}

type BatchCommittedResult struct {
	TransactionContexts []*transaction.Context
	BatchNum            int64
	BatchMetadata       []byte
	WrittenSignal       chan struct{}
	SyncedSignal        chan struct{}
}

type MemStorage struct {
	store sync.Map
}

func (ms *MemStorage) GetState(ns, key string) ([]byte, error) {
	value, _ := ms.store.Load(ns + key)
	if value == nil {
		return nil, nil
	}
	return value.([]byte), nil
}

func (ms *MemStorage) Write(batchCommittedResult *BatchCommittedResult) {
	for _, tctx := range batchCommittedResult.TransactionContexts {
		if tctx.Result.ResultCode == transaction.TxResultValid {
			for _, kvWrite := range tctx.RWSet.WSet {
				if kvWrite.IsDelete {
					ms.store.Delete(kvWrite.Key)
				} else {
					ms.store.Store(tctx.Transaction.GetContractID()+kvWrite.Key, kvWrite.Value)
				}
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
