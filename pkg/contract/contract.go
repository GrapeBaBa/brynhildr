package contract

type Response interface {
	GetStatus() int32

	GetMessage() string

	GetPayload() []byte
}

type Contract interface {
	Init(contractTransactionContext TransactionContext) Response

	Invoke(contractTransactionContext TransactionContext) Response
}

type TransactionContext interface {
	GetContractCallStub() CallStub
}

type CallStub interface {
	PutState(key string, value []byte) error

	GetState(key string) ([]byte, error)

	DelState(key string) error
}
