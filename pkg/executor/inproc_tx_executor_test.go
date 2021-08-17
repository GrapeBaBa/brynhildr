package executor

import (
	"github.com/GrapeBaBa/brynhildr/pkg/contract"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInProcTransactionExecutor(t *testing.T) {
	contracts := map[string]contract.Contract{
		"testContract": &TestContract{},
	}
	te := NewInProcTransactionExecutor(&storage.MemStorage{}, contracts)
	assert.Equal(t, len(te.contracts), 1)
}

func TestInProcTransactionExecutor_Execute(t *testing.T) {
	contracts := map[string]contract.Contract{
		"testContract": &TestContract{},
	}
	te := NewInProcTransactionExecutor(&storage.MemStorage{}, contracts)

	testTransaction := &transaction.Int64IDTransaction{
		Id:         &transaction.Int64TID{Id: 0},
		ExecType:   contract.InProcTxExec,
		ContractId: "testContract",
		Method:     "invoke",
		Args:       []string{"addState", "key", "value", "payload"},
	}
	testTransactionContext := &transaction.Context{
		Transaction: testTransaction,
		RWSet:       &transaction.RWSet{RSet: make([]transaction.KVRead, 0), WSet: make([]transaction.KVWrite, 0)},
		Result:      &transaction.Result{},
	}
	te.Execute(testTransactionContext)
	assert.Equal(t, string(testTransactionContext.Result.ExecResult.GetPayload()), "payload")
	assert.Equal(t, testTransactionContext.RWSet.WSet[0].Key, "key")
	assert.Equal(t, string(testTransactionContext.RWSet.WSet[0].Value), "value")
	assert.Equal(t, testTransactionContext.RWSet.WSet[0].IsDelete, false)
}
