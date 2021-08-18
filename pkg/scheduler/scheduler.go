package scheduler

import (
	"github.com/GrapeBaBa/brynhildr/pkg/committer"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type Scheduler interface {
	Execute(batch transaction.Batch)
	Commit(batchExecutionResult *committer.BatchExecutionResult)
	Flush(batchCommitResult *storage.BatchCommittedResult)
}
