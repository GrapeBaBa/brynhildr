package committer

import (
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"sync"
)

type BatchExecutionResult struct {
	TransactionContexts []*transaction.Context
	BatchNum            int64
	BatchMetadata       []byte
	ReserveWritesTable  *sync.Map
}

type BatchCommitter interface {
	// Commit checks the transaction conflict and generate transaction commit status.
	Commit(batchExecutionResult *BatchExecutionResult) *storage.BatchCommittedResult
}
