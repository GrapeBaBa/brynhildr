package contract

import (
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/GrapeBaBa/brynhildr/pkg/wsetcache"
)

type InProcResponse struct {
	Status int32

	Message string

	Payload []byte
}

func (ipr *InProcResponse) GetStatus() int32 {
	return ipr.Status
}

func (ipr *InProcResponse) GetMessage() string {
	return ipr.Message
}

func (ipr *InProcResponse) GetPayload() []byte {
	return ipr.Payload
}

type InProcContractCallStub struct {
	execTranContext *transaction.Context
	writeCache      wsetcache.WriteSetCache
	storageSnapshot storage.Storage
}

func NewInProcContractCallStub(execTranContext *transaction.Context, writeCache wsetcache.WriteSetCache, storageSnapshot storage.Storage) *InProcContractCallStub {
	ipcs := &InProcContractCallStub{
		execTranContext: execTranContext,
		writeCache:      writeCache,
		storageSnapshot: storageSnapshot,
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
	kvWrite := ipcs.writeCache.GetState(key)
	if kvWrite.Key == "" {
		return ipcs.storageSnapshot.GetState(key)
	} else {
		if kvWrite.IsDelete {
			return nil, nil
		} else {
			return kvWrite.Value, nil
		}
	}
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
