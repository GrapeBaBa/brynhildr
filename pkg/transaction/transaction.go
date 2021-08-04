package transaction

const (
	ContractInProc = iota

	TxResultValid = 0
)

type KVRead struct {
	Key string
}

type KVWrite struct {
	Key      string
	IsDelete bool
	Value    []byte
}

type RWSet struct {
	RSet []*KVRead
	WSet []*KVWrite
}

type TID interface {
	CompareTo(anotherTID TID) int
}

type Transaction interface {
	GetTID() TID
	GetExecutorType() int
	GetContractID() string
	GetMethod() string
}

type Batch interface {
	GetTransactions() []Transaction
}

type Context struct {
	TX         Transaction
	RWSet      *RWSet
	ResultCode int32
}
