package committer

import (
	"context"
	"github.com/GrapeBaBa/brynhild/pkg/storage"
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
	"github.com/GrapeBaBa/brynhild/pkg/wsetcache"
	"sync"
	"sync/atomic"
)

type Flusher struct {
	writeSetCacheQueue *CopyOnWriteBatchAndStateQueue
	storage            storage.Storage
}

func (flusher *Flusher) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		size := flusher.writeSetCacheQueue.size()
		batchAndStates := flusher.writeSetCacheQueue.getFirstN(size)
		flusher.storage.Commit(batchAndStates)
		flusher.writeSetCacheQueue.removeFirstN(size)
	}

}

type ParallelTransactionCommitter struct {
	batchAndStateQueue *CopyOnWriteBatchAndStateQueue
	reserveWriteTable  *sync.Map
	storage            storage.Storage
	buildWriteSetCache func() wsetcache.WriteSetCache
}

type CopyOnWriteBatchAndStateQueue struct {
	queue atomic.Value
	mutex sync.Mutex
}

func NewCopyOnWriteWriteSetCacheQueue() *CopyOnWriteBatchAndStateQueue {
	queueSlice := make([]storage.BatchAndUpdatedState, 0)
	cowq := &CopyOnWriteBatchAndStateQueue{}
	cowq.queue.Store(queueSlice)

	return cowq
}

func (cowq *CopyOnWriteBatchAndStateQueue) addLast(baus storage.BatchAndUpdatedState) {
	cowq.mutex.Lock()
	defer cowq.mutex.Unlock()
	oldSlice := cowq.queue.Load().([]storage.BatchAndUpdatedState)
	oldLen := len(oldSlice)
	newSlice := make([]storage.BatchAndUpdatedState, oldLen+1)
	copy(newSlice, oldSlice)
	newSlice[oldLen] = baus
	cowq.queue.Store(newSlice)
}

func (cowq *CopyOnWriteBatchAndStateQueue) getFirstN(len int) []storage.BatchAndUpdatedState {
	slice := cowq.queue.Load().([]storage.BatchAndUpdatedState)
	return slice[0 : len-1]
}

func (cowq *CopyOnWriteBatchAndStateQueue) getFirst() storage.BatchAndUpdatedState {
	slice := cowq.queue.Load().([]storage.BatchAndUpdatedState)
	return slice[0]
}

func (cowq *CopyOnWriteBatchAndStateQueue) getLast() storage.BatchAndUpdatedState {
	slice := cowq.queue.Load().([]storage.BatchAndUpdatedState)
	return slice[len(slice)-1]
}

func (cowq *CopyOnWriteBatchAndStateQueue) removeFirst() {
	cowq.mutex.Lock()
	defer cowq.mutex.Unlock()
	oldSlice := cowq.queue.Load().([]storage.BatchAndUpdatedState)
	oldLen := len(oldSlice)
	newSlice := make([]storage.BatchAndUpdatedState, oldLen-1)
	copy(newSlice, oldSlice[1:])
	cowq.queue.Store(newSlice)
}

func (cowq *CopyOnWriteBatchAndStateQueue) removeFirstN(size int) {
	cowq.mutex.Lock()
	defer cowq.mutex.Unlock()
	oldSlice := cowq.queue.Load().([]storage.BatchAndUpdatedState)
	oldLen := len(oldSlice)
	newSlice := make([]storage.BatchAndUpdatedState, oldLen-size)
	if oldLen > size {
		copy(newSlice, oldSlice[size:])
	}

	cowq.queue.Store(newSlice)
}

func (cowq *CopyOnWriteBatchAndStateQueue) size() int {
	slice := cowq.queue.Load().([]storage.BatchAndUpdatedState)
	return len(slice)
}

func (ptc *ParallelTransactionCommitter) Commit(ctxs []*transaction.Context) {
	writeSetCache := ptc.buildWriteSetCache()
	var wg sync.WaitGroup
	for _, ctx := range ctxs {
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
			exist := hasConflict(keysSlice, ctx.TX.GetTID(), ptc.reserveWriteTable)
			if !exist {
				ctx.ResultCode = transaction.TxResultValid
				for _, kvWrite := range ctx.RWSet.WSet {
					writeSetCache.PutState(kvWrite.Key, *kvWrite)
				}
			} else {

			}
			wg.Done()
		}(ctx, &wg)
	}

	wg.Wait()
	batchAndState := storage.BatchAndUpdatedState{TransactionContexts: ctxs, UpdatedState: writeSetCache}
	ptc.batchAndStateQueue.addLast(batchAndState)
}

func hasConflict(keys []string, tid transaction.TID, reserveWriteTable *sync.Map) bool {
	for _, key := range keys {
		if reservedTID, ok := reserveWriteTable.Load(key); ok && reservedTID.(*atomic.Value).Load().(transaction.TID).CompareTo(tid) < 0 {
			return true
		}
	}
	return false
}
