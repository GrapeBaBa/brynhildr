package contract

import (
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
)

type InProcResult struct {
	Status int32

	Message string

	Payload []byte
}

func (ipr *InProcResult) GetStatus() int32 {
	return ipr.Status
}

func (ipr *InProcResult) GetMessage() string {
	return ipr.Message
}

func (ipr *InProcResult) GetPayload() []byte {
	return ipr.Payload
}

type InProcContractCallStub struct {
	execTranContext *transaction.Context
	storage         storage.Storage
}

func NewInProcContractCallStub(execTranContext *transaction.Context, storage storage.Storage) *InProcContractCallStub {
	ipcs := &InProcContractCallStub{
		execTranContext: execTranContext,
		storage:         storage,
	}

	return ipcs
}

func (ipcs *InProcContractCallStub) PutState(key string, value []byte) error {
	kvWrite := transaction.KVWrite{Key: key, IsDelete: false, Value: value}
	ipcs.execTranContext.RWSet.WSet = append(ipcs.execTranContext.RWSet.WSet, kvWrite)
	return nil
}

func (ipcs *InProcContractCallStub) DelState(key string) error {
	kvWrite := transaction.KVWrite{Key: key, IsDelete: true, Value: nil}
	ipcs.execTranContext.RWSet.WSet = append(ipcs.execTranContext.RWSet.WSet, kvWrite)
	return nil
}

func (ipcs *InProcContractCallStub) GetState(key string) ([]byte, error) {
	return ipcs.storage.GetState(key)
}

func (ipcs *InProcContractCallStub) GetFunctionAndArgs() (string, []string) {
	return ipcs.execTranContext.Transaction.GetFunctionAndArgs()
}

type InProcContractTransactionContext struct {
	inProcContractCallStub *InProcContractCallStub
}

func NewInProcContractTransactionContext(inProcContractCallStub *InProcContractCallStub) *InProcContractTransactionContext {
	ipctc := &InProcContractTransactionContext{
		inProcContractCallStub: inProcContractCallStub,
	}

	return ipctc
}

func (ipctc *InProcContractTransactionContext) GetContractCallStub() CallStub {
	return ipctc.inProcContractCallStub
}
