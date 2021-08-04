package storage

import (
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
	"github.com/GrapeBaBa/brynhild/pkg/wsetcache"
)

type BatchAndUpdatedState struct {
	TransactionContexts []*transaction.Context
	UpdatedState        wsetcache.WriteSetCache
}

type Storage interface {
	GetState(key string) ([]byte, error)

	Commit(batchAndStates []BatchAndUpdatedState)
}
