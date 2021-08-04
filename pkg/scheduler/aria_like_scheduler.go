package scheduler

import (
	"github.com/GrapeBaBa/brynhild/pkg/committer"
	"github.com/GrapeBaBa/brynhild/pkg/executor"
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
	"github.com/pingcap/tidb/util/bitmap"
	"golang.org/x/sync/semaphore"
	"sync"
	"sync/atomic"
)

type AriaLikeScheduler struct {
	reserveWriteTable *sync.Map
	semp              *semaphore.Weighted
	executorMgr       executor.Manager
	batchBitMap       *bitmap.ConcurrentBitmap
	committer         committer.Committer
}

func NewAriaLikeScheduler(concurLimit int64) *AriaLikeScheduler {
	as := &AriaLikeScheduler{
		semp: semaphore.NewWeighted(concurLimit),
	}

	return as
}

func (as *AriaLikeScheduler) reserveWrites(ctx *transaction.Context) {
	ctxTID := ctx.TX.GetTID()
	for _, write := range ctx.RWSet.WSet {
		var currTIDValue atomic.Value
		currTIDValue.Store(ctxTID)
		// First store current tid for write key, it will success when this key is not exist previous
		existTIDValue, loaded := as.reserveWriteTable.LoadOrStore(write.Key, currTIDValue)
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
						existTIDValue, _ = as.reserveWriteTable.Load(write.Key)
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

func (as *AriaLikeScheduler) execute(ctxs []*transaction.Context) {
	var wg sync.WaitGroup
	for _, tx := range ctxs {
		wg.Add(1)
		go func(ctx *transaction.Context, wg *sync.WaitGroup) {
			as.executorMgr.Execute(ctx)
			as.reserveWrites(ctx)
			wg.Done()
		}(tx, &wg)
	}

	wg.Wait()
}

func (as *AriaLikeScheduler) commit(ctxs []*transaction.Context) {
	as.committer.Commit(ctxs)
}