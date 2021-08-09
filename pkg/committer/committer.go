package committer

import (
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type BatchCommitter interface {
	Commit(batchAndWSet *transaction.BatchAndWSet)
}
