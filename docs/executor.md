# Batch Executor

```go
type BatchExecutor interface {
// Execute executes the a batch of transactions and generate the execution result(rwset)
Execute(batch transaction.Batch) *transaction.BatchAndWSet
}
```

Batch execution phase is the first phase in Brynhildr, it receives a transaction batch and execute all the transactions.

## AriaLikeBatchExecutor
```go
type AriaLikeBatchExecutor struct {
	txExecMgr         *TransactionExecutorManager
	reserveWriteTable *sync.Map
}
```
Refer to [Aria: A Fast and Practical Deterministic OLTP Database](http://www.vldb.org/pvldb/vol13/p2047-lu.pdf), this executor will execute all the transaction parallel without DAG dependency analysis.
Each transaction will generate read-write set after execution, once all transaction execution finished, it will generate a global reserved write set table. The reserved write set table store the all write 
set key with the smallest batch number.

# Transaction Executor

```go
type TransactionExecutor interface {
// Execute executes a transaction and generate the execution result(rwset)
Execute(context *transaction.Context)
}
```

Transaction executor is responsible for specific transaction execution. It can be implemented by any contract engine
technology.


