package transaction

const (
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
	GetFunctionAndArgs() (string, []string)
}

type Batch interface {
	GetTransactions() []Transaction
	GetNumber() int64
	GetMetadata() []byte
}

type ExecResult interface {
	GetStatus() int32

	GetMessage() string

	GetPayload() []byte
}

type Context struct {
	Transaction Transaction
	RWSet       *RWSet
	Result      *Result
}

type Result struct {
	ResultCode int32
	ExecResult ExecResult
}

type Int64TID struct {
	Id int64
}

func (tid *Int64TID) CompareTo(anotherTID TID) int {
	if tid.Id > anotherTID.(*Int64TID).Id {
		return 1
	} else if tid.Id < anotherTID.(*Int64TID).Id {
		return -1
	} else {
		return 0
	}
}

type Int64IDTransaction struct {
	Id         *Int64TID
	ExecType   int
	ContractId string
	Method     string
	Args       []string
}

func (tt *Int64IDTransaction) GetTID() TID {
	return tt.Id
}

func (tt *Int64IDTransaction) GetExecutorType() int {
	return tt.ExecType
}

func (tt *Int64IDTransaction) GetContractID() string {
	return tt.ContractId
}

func (tt *Int64IDTransaction) GetMethod() string {
	return tt.Method
}

func (tt *Int64IDTransaction) GetFunctionAndArgs() (string, []string) {
	return tt.Args[0], tt.Args[1:]
}

type Int64Batch struct {
	Transactions []*Int64IDTransaction
	Number       int64
}

func (ib *Int64Batch) GetTransactions() []Transaction {
	txs := make([]Transaction, 0)
	for _, tx := range ib.Transactions {
		txs = append(txs, tx)
	}

	return txs
}

func (ib *Int64Batch) GetNumber() int64 {
	return ib.Number
}

func (ib *Int64Batch) GetMetadata() []byte {
	return nil
}
