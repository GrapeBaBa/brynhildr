package scheduler

import (
	"context"
	"github.com/GrapeBaBa/brynhildr/pkg/contract"
	"github.com/GrapeBaBa/brynhildr/pkg/executor"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func TestNewAriaLikeScheduler(t *testing.T) {
	contracts := map[string]contract.Contract{
		"testContract": &TestContract{},
	}
	executors := map[int]executor.TransactionExecutor{
		contract.InProcTxExec: executor.NewInProcTransactionExecutor(&storage.MemStorage{}, contracts),
	}

	tem := executor.NewTransactionExecutorManager(executors)
	store := &storage.MemStorage{}
	scheduler := NewAriaLikeScheduler(tem, store)
	assert.Exactly(t, scheduler.storage, store)
}

func TestAriaLikeScheduler_Start(t *testing.T) {
	contracts := map[string]contract.Contract{
		"testContract": &TestContract{},
	}
	executors := map[int]executor.TransactionExecutor{
		contract.InProcTxExec: executor.NewInProcTransactionExecutor(&storage.MemStorage{}, contracts),
	}

	tem := executor.NewTransactionExecutorManager(executors)
	store := &storage.MemStorage{}
	scheduler := NewAriaLikeScheduler(tem, store)
	ctx, canFunc := context.WithCancel(context.Background())
	scheduler.Start(ctx)

	go func() {
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
			Args:       []string{"addState", "key", "value1", "payload1"},
		}}}
		scheduler.Handle(batch)
	}()

	time.Sleep(1 * time.Second)
	value, _ := store.GetState("testContract", "key")
	assert.Equal(t, string(value), "value")
	canFunc()
}
