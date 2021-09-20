package executor

import (
	"sync"
	"sync/atomic"

	"github.com/GrapeBaBa/brynhildr/pkg/committer"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type AriaLikeBatchExecutor struct {
	txExecMgr *TransactionExecutorManager
}

func NewAriaLikeBatchExecutor(txExecMgr *TransactionExecutorManager) *AriaLikeBatchExecutor {
	return &AriaLikeBatchExecutor{
		txExecMgr: txExecMgr,
	}
}

func (abe *AriaLikeBatchExecutor) Execute(batch transaction.Batch) *committer.BatchExecutionResult {
	reserveWriteTable := &sync.Map{}
	var wg sync.WaitGroup
	txs := batch.GetTransactions()
	tctxs := make([]*transaction.Context, len(txs))
	for i, tx := range txs {
		tctx := &transaction.Context{Transaction: tx, RWSet: &transaction.RWSet{RSet: make([]transaction.KVRead, 0), WSet: make([]transaction.KVWrite, 0)}, Result: &transaction.Result{}}
		tctxs[i] = tctx
		wg.Add(1)
		go func(ctx *transaction.Context, wg *sync.WaitGroup) {
			abe.txExecMgr.Execute(ctx)
			reserveWrites(ctx, reserveWriteTable)
			wg.Done()
		}(tctx, &wg)
	}

	wg.Wait()

	immutableReserveWriteTable := &sync.Map{}
	reserveWriteTable.Range(func(key, value interface{}) bool {
		immutableReserveWriteTable.Store(key, value)
		return true
	})
	batchAndUpdatedState := &committer.BatchExecutionResult{TransactionContexts: tctxs, BatchNum: batch.GetNumber(), BatchMetadata: batch.GetMetadata(), ReserveWritesTable: immutableReserveWriteTable}
	return batchAndUpdatedState
}

func reserveWrites(ctx *transaction.Context, reserveWriteTable *sync.Map) {
	ctxTID := ctx.Transaction.GetTID()
	for _, write := range ctx.RWSet.WSet {
		var currTIDValue atomic.Value
		currTIDValue.Store(ctxTID)
		// First store current tid for write key, it will success when this key is not exist previous
		existTIDValue, loaded := reserveWriteTable.LoadOrStore(write.Key, &currTIDValue)
		// This key is already exist
		if loaded {
			existTID := existTIDValue.(*atomic.Value).Load().(transaction.TID)
			// Compare current tid and existed tid, if current tid is smaller
			if ctxTID.CompareTo(existTID) < 0 {
				// Atomic store current tid in reserveWriteTable
				swapped := existTIDValue.(*atomic.Value).CompareAndSwap(existTID, ctxTID)
				// If not store success
				if !swapped {
					// Try to store current tid if still needed
					for {
						existTIDValue, _ = reserveWriteTable.Load(write.Key)
						existTID = existTIDValue.(*atomic.Value).Load().(transaction.TID)
						if ctxTID.CompareTo(existTID) < 0 {
							swapped = existTIDValue.(*atomic.Value).CompareAndSwap(existTID, ctxTID)
							if swapped {
								break
							}
						} else {
							break
						}
					}
				}
			}
		}
	}
}
