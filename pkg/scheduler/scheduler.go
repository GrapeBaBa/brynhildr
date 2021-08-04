package scheduler

import (
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
)

type Scheduler interface {
	Execute(batch transaction.Batch)
	Commit(batch transaction.Batch)
}
