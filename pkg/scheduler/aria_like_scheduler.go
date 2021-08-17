package scheduler

import (
	"context"
	"github.com/GrapeBaBa/brynhildr/pkg/committer"
	"github.com/GrapeBaBa/brynhildr/pkg/executor"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"sync"
)

type AriaLikeScheduler struct {
	waitToExecuteCh chan transaction.Batch
	waitToCommitCh  chan *committer.BatchExecutionResult
	waitToFlushCh   chan *storage.BatchCommittedResult
	readyToExecCh   chan struct{}
	batchExecutor   executor.BatchExecutor
	committer       committer.BatchCommitter
	storage         storage.Storage
}

func NewAriaLikeScheduler(txExecMgr *executor.TransactionExecutorManager, store storage.Storage) *AriaLikeScheduler {
	reserveWriteTable := &sync.Map{}
	batchExecutor := executor.NewAriaLikeBatchExecutor(txExecMgr, reserveWriteTable)
	comm := committer.NewAriaLikeBatchCommitter(reserveWriteTable)
	as := &AriaLikeScheduler{
		waitToFlushCh:   make(chan *storage.BatchCommittedResult),
		waitToCommitCh:  make(chan *committer.BatchExecutionResult),
		waitToExecuteCh: make(chan transaction.Batch),
		readyToExecCh:   make(chan struct{}),
		batchExecutor:   batchExecutor,
		committer:       comm,
		storage:         store,
	}

	return as
}

func (as *AriaLikeScheduler) Handle(batch transaction.Batch) {
	as.waitToExecuteCh <- batch
}

func (as *AriaLikeScheduler) Start(ctx context.Context) {
	go func() {
		as.readyToExecCh <- struct{}{}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case batch := <-as.waitToExecuteCh:
				<-as.readyToExecCh
				batchExecRes := as.batchExecutor.Execute(batch)
				as.waitToCommitCh <- batchExecRes
			}
		}

	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case batchExecRes := <-as.waitToCommitCh:
				batchCommitRes := as.committer.Commit(batchExecRes)
				batchCommitRes.WrittenSignal = as.readyToExecCh
				as.waitToFlushCh <- batchCommitRes
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case batchCommitRes := <-as.waitToFlushCh:
				as.storage.Write(batchCommitRes)
			}
		}
	}()

}
