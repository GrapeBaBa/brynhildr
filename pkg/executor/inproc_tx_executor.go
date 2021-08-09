package executor

import (
	"github.com/GrapeBaBa/brynhildr/pkg/contract"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"github.com/GrapeBaBa/brynhildr/pkg/wsetcache"
	"strings"
)

type InProcTransactionExecutor struct {
	writeSetCache   wsetcache.WriteSetCache
	storageSnapshot storage.Storage
	contracts       map[string]contract.Contract
}

func (ipe *InProcTransactionExecutor) Execute(context *transaction.Context) {
	tx := context.TX
	//TODO:Prevent contract panic
	ipccs := contract.NewInProcContractCallStub(context, ipe.writeSetCache, ipe.storageSnapshot)
	ctc := contract.NewInProcContractTransactionContext(ipccs)
	// no use reflection for performance consideration
	if strings.EqualFold(tx.GetMethod(), "init") {
		ipe.contracts[tx.GetContractID()].Init(ctc)
	} else {
		ipe.contracts[tx.GetContractID()].Invoke(ctc)
	}
}

func (ae *TransactionExecutorManager) Execute(context *transaction.Context) {
	ae.executors[context.TX.GetExecutorType()].Execute(context)
}
