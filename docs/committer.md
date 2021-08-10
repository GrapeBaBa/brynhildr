# Batch Committer

```go
type BatchCommitter interface {
// Commit checks the transaction conflict and generate transaction commit status.
Commit(batchAndWSet *BatchAndWSet)
}
```

Batch commit phase is the second phase in Brynhildr, it receives a transaction result batch and generates transaction commitment status.

## AriaLikeBatchExecutor

```go
type AriaLikeBatchCommitter struct {
reserveWriteTable  *sync.Map
waitToWriteCh      chan BatchAndWSet
buildWriteSetCache func() wsetcache.WriteSetCache
}
```

Refer to [Aria: A Fast and Practical Deterministic OLTP Database](http://www.vldb.org/pvldb/vol13/p2047-lu.pdf), this
committer will commit all the transaction parallel. Each transaction will generate commitment status, once all transactions' commitment finished, it will generate a global write set
cache for this batch which includes all successful commitment transaction write set.

