package executor

import (
	"github.com/GrapeBaBa/brynhildr/pkg/committer"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type BatchExecutor interface {
	// Execute executes the a batch of transactions and generate the execution result(rwset)
	Execute(batch transaction.Batch) *committer.BatchExecutionResult
}

type TransactionExecutor interface {
	// Execute executes a transaction and generate the execution result(rwset)
	Execute(context *transaction.Context)
}

type TransactionExecutorManager struct {
	executors map[int]TransactionExecutor
}

func (tem *TransactionExecutorManager) Execute(context *transaction.Context) {
	tem.executors[context.Transaction.GetExecutorType()].Execute(context)
}

func NewTransactionExecutorManager(executors map[int]TransactionExecutor) *TransactionExecutorManager {
	return &TransactionExecutorManager{
		executors: executors,
	}
}
