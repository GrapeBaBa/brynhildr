package executor

import (
	"github.com/GrapeBaBa/brynhildr/pkg/contract"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestContract struct {
}

func (tc *TestContract) Init(contractTransactionContext contract.TransactionContext) transaction.ExecResult {
	return nil
}

func (tc *TestContract) Invoke(contractTransactionContext contract.TransactionContext) transaction.ExecResult {
	f, args := contractTransactionContext.GetContractCallStub().GetFunctionAndArgs()
	if f == "addState" {
		return tc.addState(args, contractTransactionContext.GetContractCallStub())
	}

	return nil
}

func (tc *TestContract) addState(args []string, stub contract.CallStub) transaction.ExecResult {
	_ = stub.PutState(args[0], []byte(args[1]))
	return &contract.InProcResult{Status: contract.ExecStatusSuccess, Payload: []byte(args[2])}
}

func TestNewTransactionExecutorManager(t *testing.T) {
	contracts := map[string]contract.Contract{
		"testContract": &TestContract{},
	}
	executors := map[int]TransactionExecutor{
		contract.InProcTxExec: NewInProcTransactionExecutor(&storage.MemStorage{}, contracts),
	}

	tem := NewTransactionExecutorManager(executors)
	assert.Equal(t, len(tem.executors), 1)
}

func TestTransactionExecutorManager_Execute(t *testing.T) {
	contracts := map[string]contract.Contract{
		"testContract": &TestContract{},
	}
	executors := map[int]TransactionExecutor{
		contract.InProcTxExec: NewInProcTransactionExecutor(&storage.MemStorage{}, contracts),
	}

	tem := NewTransactionExecutorManager(executors)

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
	tem.Execute(testTransactionContext)
	assert.Equal(t, string(testTransactionContext.Result.ExecResult.GetPayload()), "payload")
	assert.Equal(t, testTransactionContext.RWSet.WSet[0].Key, "key")
	assert.Equal(t, string(testTransactionContext.RWSet.WSet[0].Value), "value")
	assert.Equal(t, testTransactionContext.RWSet.WSet[0].IsDelete, false)
}
