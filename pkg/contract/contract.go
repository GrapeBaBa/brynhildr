package contract

import "github.com/GrapeBaBa/brynhildr/pkg/transaction"

const (
	ExecStatusSuccess = iota
)

const (
	InProcTxExec = iota
)

type Contract interface {
	Init(contractTransactionContext TransactionContext) transaction.ExecResult

	Invoke(contractTransactionContext TransactionContext) transaction.ExecResult
}

type TransactionContext interface {
	GetContractCallStub() CallStub
}

type CallStub interface {
	PutState(key string, value []byte) error

	GetState(key string) ([]byte, error)

	DelState(key string) error

	GetFunctionAndArgs() (string, []string)
}
