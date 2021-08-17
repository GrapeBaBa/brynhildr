package executor

import (
	"github.com/GrapeBaBa/brynhildr/pkg/contract"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewAriaLikeBatchExecutor(t *testing.T) {
	contracts := map[string]contract.Contract{
		"testContract": &TestContract{},
	}
	executors := map[int]TransactionExecutor{
		contract.InProcTxExec: NewInProcTransactionExecutor(&storage.MemStorage{}, contracts),
	}

	tem := NewTransactionExecutorManager(executors)
	mem := &sync.Map{}
	albe := NewAriaLikeBatchExecutor(tem, mem)
	batch := &transaction.Int64Batch{Number: 10, Transactions: []*transaction.Int64IDTransaction{{
		Id:         &transaction.Int64TID{Id: 100},
		ExecType:   contract.InProcTxExec,
		ContractId: "testContract",
		Method:     "invoke",
		Args:       []string{"addState", "key", "value", "payload"},
	}, {
		Id:         &transaction.Int64TID{Id: 101},
		ExecType:   contract.InProcTxExec,
		ContractId: "testContract",
		Method:     "invoke",
		Args:       []string{"addState", "key", "value", "payload1"},
	}}}
	res := albe.Execute(batch)
	assert.Equal(t, int64(10), res.BatchNum)
	confKey, _ := albe.reserveWriteTable.Load("key")
	tid := confKey.(*atomic.Value).Load().(*transaction.Int64TID)
	assert.Equal(t, int64(100), tid.Id)

}
