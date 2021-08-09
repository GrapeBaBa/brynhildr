package scheduler

import (
	"context"
	"github.com/GrapeBaBa/brynhildr/pkg/committer"
	"github.com/GrapeBaBa/brynhildr/pkg/executor"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/pingcap/tidb/util/bitmap"
	"golang.org/x/sync/semaphore"
	"sync"
)

type AriaLikeScheduler struct {
	waitToExecuteCh   chan transaction.Batch
	waitToCommitCh    chan *committer.BatchAndWSet
	waitToFlushCh     chan *storage.BatchAndWSetSyncer
	readyToExecCh     chan struct{}
	reserveWriteTable *sync.Map
	semp              *semaphore.Weighted
	batchExecutor     executor.BatchExecutor
	batchBitMap       *bitmap.ConcurrentBitmap
	committer         committer.BatchCommitter
	storage           storage.Storage
}

func NewAriaLikeScheduler(concurLimit int64) *AriaLikeScheduler {
	as := &AriaLikeScheduler{
		semp: semaphore.NewWeighted(concurLimit),
	}

	return as
}

func (as *AriaLikeScheduler) Receive(batch transaction.Batch) {
	as.waitToExecuteCh <- batch
}

func (as *AriaLikeScheduler) NotifyExec() {
	as.readyToExecCh <- struct{}{}
}

func (as *AriaLikeScheduler) Start(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case batch := <-as.waitToExecuteCh:
				<-as.readyToExecCh
				batchAndWSet := as.batchExecutor.Execute(batch)
				as.waitToCommitCh <- batchAndWSet
			}
		}

	}(ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case batchAndWSet := <-as.waitToCommitCh:
				as.committer.Commit(batchAndWSet)
				syncer := &storage.BatchAndWSetSyncer{BatchAndWSet: *batchAndWSet, WrittenSignal: as.readyToExecCh}
				as.waitToFlushCh <- syncer
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case syncer := <-as.waitToFlushCh:
				as.storage.Write(syncer)
			}
		}
	}(ctx)
}
