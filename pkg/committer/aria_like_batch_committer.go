package committer

import (
	"sync"
	"sync/atomic"

	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type AriaLikeBatchCommitter struct {
	waitToWriteCh chan BatchExecutionResult
}

func NewAriaLikeBatchCommitter() *AriaLikeBatchCommitter {
	return &AriaLikeBatchCommitter{
		waitToWriteCh: make(chan BatchExecutionResult),
	}
}

func (ptc *AriaLikeBatchCommitter) Commit(batchExecutionResult *BatchExecutionResult) *storage.BatchCommittedResult {
	var wg sync.WaitGroup
	res := &storage.BatchCommittedResult{BatchNum: batchExecutionResult.BatchNum, BatchMetadata: batchExecutionResult.BatchMetadata, TransactionContexts: make([]*transaction.Context, len(batchExecutionResult.TransactionContexts))}
	for i, ctx := range batchExecutionResult.TransactionContexts {
		res.TransactionContexts[i] = ctx
		wg.Add(1)
		go func(ctx *transaction.Context, wg *sync.WaitGroup) {
			keysSlice := make([]string, 0)
			keysMap := make(map[string]int8)
			for _, kvRead := range ctx.RWSet.RSet {
				if _, ok := keysMap[kvRead.Key]; !ok {
					keysMap[kvRead.Key] = 1
					keysSlice = append(keysSlice, kvRead.Key)
				}
			}

			for _, kvWrite := range ctx.RWSet.WSet {
				if _, ok := keysMap[kvWrite.Key]; !ok {
					keysMap[kvWrite.Key] = 1
					keysSlice = append(keysSlice, kvWrite.Key)
				}
			}
			keysMap = nil
			exist := hasConflict(keysSlice, ctx.Transaction.GetTID(), batchExecutionResult.ReserveWritesTable)
			if !exist {
				ctx.Result.ResultCode = transaction.TxResultValid
			} else {
				ctx.Result.ResultCode = transaction.TxResultDependencyConflict
			}
			wg.Done()
		}(ctx, &wg)
	}

	wg.Wait()
	return res
}

func hasConflict(keys []string, tid transaction.TID, reserveWriteTable *sync.Map) bool {
	for _, key := range keys {
		if reservedTID, ok := reserveWriteTable.Load(key); ok && reservedTID.(*atomic.Value).Load().(transaction.TID).CompareTo(tid) < 0 {
			return true
		}
	}
	return false
}
