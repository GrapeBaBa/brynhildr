package scheduler

import (
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type Scheduler interface {
	Execute(batch transaction.Batch)
	Commit(batchAndWSet *transaction.BatchAndWSet)
}
