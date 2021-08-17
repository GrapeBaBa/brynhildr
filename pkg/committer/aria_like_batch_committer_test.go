package committer

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/GrapeBaBa/brynhildr/pkg/contract"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/stretchr/testify/assert"
)

func TestNewAriaLikeBatchCommitter(t *testing.T) {
	mem := &sync.Map{}
	committer := NewAriaLikeBatchCommitter(mem)
	assert.Exactly(t, committer.reserveWriteTable, mem)
	assert.Equal(t, len(committer.waitToWriteCh), 0)
}

func TestAriaLikeBatchCommitter_Commit(t *testing.T) {
	var resKey atomic.Value
	resKey.Store(&transaction.Int64TID{Id: 100})
	mem := &sync.Map{}
	mem.Store("key", &resKey)

	beres := &BatchExecutionResult{
		BatchNum: 20,
		TransactionContexts: []*transaction.Context{
			{
				Transaction: &transaction.Int64IDTransaction{
					Id:         &transaction.Int64TID{Id: 100},
					ExecType:   contract.InProcTxExec,
					ContractId: "testContract",
					Method:     "invoke",
					Args:       []string{"addState", "key", "value", "payload"},
				},
				RWSet: &transaction.RWSet{
					RSet: []transaction.KVRead{},
					WSet: []transaction.KVWrite{
						{
							Key:      "key",
							IsDelete: false,
							Value:    []byte("value"),
						},
					},
				},
				Result: &transaction.Result{
					ExecResult: &contract.InProcResult{
						Payload: []byte("payload"),
						Status:  contract.ExecStatusSuccess,
					},
				},
			},
			{
				Transaction: &transaction.Int64IDTransaction{
					Id:         &transaction.Int64TID{Id: 101},
					ExecType:   contract.InProcTxExec,
					ContractId: "testContract",
					Method:     "invoke",
					Args:       []string{"addState", "key", "value", "payload"},
				},
				RWSet: &transaction.RWSet{
					RSet: []transaction.KVRead{},
					WSet: []transaction.KVWrite{
						{
							Key:      "key",
							IsDelete: false,
							Value:    []byte("value1"),
						},
					},
				},
				Result: &transaction.Result{
					ExecResult: &contract.InProcResult{
						Payload: []byte("payload1"),
						Status:  contract.ExecStatusSuccess,
					},
				},
			},
		},
	}
	albc := NewAriaLikeBatchCommitter(mem)
	comRes := albc.Commit(beres)
	assert.Equal(t, int64(20), comRes.BatchNum)
	assert.Equal(t, comRes.TransactionContexts[0].Result.ResultCode, int32(transaction.TxResultValid))
	assert.Equal(t, comRes.TransactionContexts[1].Result.ResultCode, int32(transaction.TxResultDependencyConflict))
}
