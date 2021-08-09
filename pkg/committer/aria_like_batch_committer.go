package committer

import (
	"sync"
	"sync/atomic"

	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/GrapeBaBa/brynhildr/pkg/wsetcache"
)

type AriaLikeBatchCommitter struct {
	reserveWriteTable  *sync.Map
	waitToWriteCh      chan BatchAndWSet
	buildWriteSetCache func() wsetcache.WriteSetCache
}

//type CopyOnWriteBatchAndStateQueue struct {
//	queue atomic.Value
//	mutex sync.Mutex
//}
//
//func NewCopyOnWriteWriteSetCacheQueue() *CopyOnWriteBatchAndStateQueue {
//	queueSlice := make([]transaction.BatchAndWSet, 0)
//	cowq := &CopyOnWriteBatchAndStateQueue{}
//	cowq.queue.Store(queueSlice)
//
//	return cowq
//}
//
//func (cowq *CopyOnWriteBatchAndStateQueue) addLast(baus transaction.BatchAndWSet) {
//	cowq.mutex.Lock()
//	defer cowq.mutex.Unlock()
//	oldSlice := cowq.queue.Load().([]transaction.BatchAndWSet)
//	oldLen := len(oldSlice)
//	newSlice := make([]transaction.BatchAndWSet, oldLen+1)
//	copy(newSlice, oldSlice)
//	newSlice[oldLen] = baus
//	cowq.queue.Store(newSlice)
//}
//
//func (cowq *CopyOnWriteBatchAndStateQueue) getFirstN(len int) []transaction.BatchAndWSet {
//	slice := cowq.queue.Load().([]transaction.BatchAndWSet)
//	return slice[0 : len-1]
//}
//
//func (cowq *CopyOnWriteBatchAndStateQueue) getFirst() transaction.BatchAndWSet {
//	slice := cowq.queue.Load().([]transaction.BatchAndWSet)
//	return slice[0]
//}
//
//func (cowq *CopyOnWriteBatchAndStateQueue) getLast() transaction.BatchAndWSet {
//	slice := cowq.queue.Load().([]transaction.BatchAndWSet)
//	return slice[len(slice)-1]
//}
//
//func (cowq *CopyOnWriteBatchAndStateQueue) removeFirst() {
//	cowq.mutex.Lock()
//	defer cowq.mutex.Unlock()
//	oldSlice := cowq.queue.Load().([]transaction.BatchAndWSet)
//	oldLen := len(oldSlice)
//	newSlice := make([]transaction.BatchAndWSet, oldLen-1)
//	copy(newSlice, oldSlice[1:])
//	cowq.queue.Store(newSlice)
//}
//
//func (cowq *CopyOnWriteBatchAndStateQueue) removeFirstN(size int) {
//	cowq.mutex.Lock()
//	defer cowq.mutex.Unlock()
//	oldSlice := cowq.queue.Load().([]transaction.BatchAndWSet)
//	oldLen := len(oldSlice)
//	newSlice := make([]transaction.BatchAndWSet, oldLen-size)
//	if oldLen > size {
//		copy(newSlice, oldSlice[size:])
//	}
//
//	cowq.queue.Store(newSlice)
//}
//
//func (cowq *CopyOnWriteBatchAndStateQueue) size() int {
//	slice := cowq.queue.Load().([]transaction.BatchAndWSet)
//	return len(slice)
//}

func (ptc *AriaLikeBatchCommitter) Commit(batchAndWSet BatchAndWSet) {
	var wg sync.WaitGroup
	for _, ctx := range batchAndWSet.TransactionContexts {
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
				ctx.Result.ResultCode = transaction.TxResultValid
				for _, kvWrite := range ctx.RWSet.WSet {
					batchAndWSet.KvWrites.PutState(kvWrite.Key, kvWrite)
				}
			} else {
				ctx.Result.ResultCode = transaction.TxResultDependencyConflict
			}
			wg.Done()
		}(ctx, &wg)
	}

	wg.Wait()
}

func hasConflict(keys []string, tid transaction.TID, reserveWriteTable *sync.Map) bool {
	for _, key := range keys {
		if reservedTID, ok := reserveWriteTable.Load(key); ok && reservedTID.(*atomic.Value).Load().(transaction.TID).CompareTo(tid) < 0 {
			return true
		}
	}
	return false
}
