package storage

import (
	"context"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/GrapeBaBa/brynhildr/pkg/wsetcache"
)

type CachedStorage struct {
	waitToWriteBatchAndWSetCh chan *transaction.BatchAndWSetSyncer
	latestWSetCache           [2]wsetcache.WriteSetCache
	immutableWSetCache        wsetcache.WriteSetCache
	storage                   Storage
}

func (cs *CachedStorage) Start(ctx context.Context) {
	//for {
	//	select {
	//	case <-ctx.Done():
	//		return
	//	case syncer := <-cs.waitToWriteBatchAndWSetCh:
	//		cs.storage.Write()
	//	}
	//}
}

func (cs *CachedStorage) Write(syncer *transaction.BatchAndWSetSyncer) {
	cs.waitToWriteBatchAndWSetCh <- syncer
}
