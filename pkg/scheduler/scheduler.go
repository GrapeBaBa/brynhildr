package scheduler

import (
	"github.com/GrapeBaBa/brynhildr/pkg/committer"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type Scheduler interface {
	Execute(batch transaction.Batch)
	Commit(batchAndWSet *committer.BatchExecutionResult)
}
