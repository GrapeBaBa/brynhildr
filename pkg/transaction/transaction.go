package transaction

import "github.com/GrapeBaBa/brynhildr/pkg/wsetcache"

const (
	ContractInProc = iota

	TxResultValid              = 0
	TxResultDependencyConflict = 1
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
	RSet []KVRead
	WSet []KVWrite
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
	GetNumber() int64
}

type Context struct {
	TX     Transaction
	RWSet  *RWSet
	Result *Result
}

type Result struct {
	ResultCode int32
}

type BatchAndWSet struct {
	TransactionContexts []*Context
	KvWrites            wsetcache.WriteSetCache
}

type BatchAndWSetSyncer struct {
	BatchAndWSet  BatchAndWSet
	WrittenSignal chan struct{}
}
