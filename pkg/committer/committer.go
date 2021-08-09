package committer

import (
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/GrapeBaBa/brynhildr/pkg/wsetcache"
)

type BatchCommitter interface {
	Commit(batchAndWSet *BatchAndWSet)
}

type BatchAndWSet struct {
	TransactionContexts []*transaction.Context
	KvWrites            wsetcache.WriteSetCache
}

