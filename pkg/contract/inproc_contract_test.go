package contract

import (
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInProcContractCallStub_GetFunctionAndArgs(t *testing.T) {
	testTransaction := &transaction.Int64IDTransaction{
		Id:         &transaction.Int64TID{Id: 0},
		ExecType:   InProcTxExec,
		ContractId: "testContract",
		Method:     "invoke",
		Args:       []string{"addState", "key", "value", "payload"},
	}
	testTransactionContext := &transaction.Context{
		Transaction: testTransaction,
		RWSet:       &transaction.RWSet{RSet: make([]transaction.KVRead, 0), WSet: make([]transaction.KVWrite, 0)},
		Result:      &transaction.Result{},
	}
	stub := NewInProcContractCallStub(testTransactionContext, &storage.MemStorage{})
	fname, args := stub.GetFunctionAndArgs()
	assert.Equal(t, "addState", fname)
	assert.Equal(t, len(args), 3)
}

func TestInProcContractCallStub_PutState(t *testing.T) {
	testTransaction := &transaction.Int64IDTransaction{
		Id:         &transaction.Int64TID{Id: 0},
		ExecType:   InProcTxExec,
		ContractId: "testContract",
		Method:     "invoke",
		Args:       []string{"addState", "key", "value", "payload"},
	}
	testTransactionContext := &transaction.Context{
		Transaction: testTransaction,
		RWSet:       &transaction.RWSet{RSet: make([]transaction.KVRead, 0), WSet: make([]transaction.KVWrite, 0)},
		Result:      &transaction.Result{},
	}
	stub := NewInProcContractCallStub(testTransactionContext, &storage.MemStorage{})
	_ = stub.PutState(testTransaction.Args[1], []byte(testTransaction.Args[2]))
	assert.Equal(t, len(testTransactionContext.RWSet.WSet), 1)
	assert.Equal(t, testTransactionContext.RWSet.WSet[0].Key, "key")
}

func TestInProcContractCallStub_DelState(t *testing.T) {
	testTransaction := &transaction.Int64IDTransaction{
		Id:         &transaction.Int64TID{Id: 0},
		ExecType:   InProcTxExec,
		ContractId: "testContract",
		Method:     "invoke",
		Args:       []string{"addState", "key", "value", "payload"},
	}
	testTransactionContext := &transaction.Context{
		Transaction: testTransaction,
		RWSet:       &transaction.RWSet{RSet: make([]transaction.KVRead, 0), WSet: make([]transaction.KVWrite, 0)},
		Result:      &transaction.Result{},
	}
	stub := NewInProcContractCallStub(testTransactionContext, &storage.MemStorage{})
	_ = stub.DelState(testTransaction.Args[1])
	assert.Equal(t, len(testTransactionContext.RWSet.WSet), 1)
	assert.Equal(t, testTransactionContext.RWSet.WSet[0].Key, "key")
	assert.Equal(t, testTransactionContext.RWSet.WSet[0].IsDelete, true)
}

func TestInProcContractCallStub_GetState(t *testing.T) {
	testTransaction := &transaction.Int64IDTransaction{
		Id:         &transaction.Int64TID{Id: 0},
		ExecType:   InProcTxExec,
		ContractId: "testContract",
		Method:     "invoke",
		Args:       []string{"addState", "key", "value", "payload"},
	}
	testTransactionContext := &transaction.Context{
		Transaction: testTransaction,
		RWSet:       &transaction.RWSet{RSet: make([]transaction.KVRead, 0), WSet: make([]transaction.KVWrite, 0)},
		Result:      &transaction.Result{},
	}
	memStore := &storage.MemStorage{}
	stub := NewInProcContractCallStub(testTransactionContext, memStore)
	_ = stub.PutState(testTransaction.Args[1], []byte(testTransaction.Args[2]))
	memStore.Write(&storage.BatchCommittedResult{TransactionContexts: []*transaction.Context{testTransactionContext}, BatchNum: 0, WrittenSignal: nil, SyncedSignal: nil})
	value, _ := stub.GetState(testTransaction.Args[1])
	assert.Equal(t, string(value), "value")
}
