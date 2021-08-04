package executor

import (
	"github.com/GrapeBaBa/brynhild/pkg/contract"
	"github.com/GrapeBaBa/brynhild/pkg/storage"
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
	"github.com/GrapeBaBa/brynhild/pkg/wsetcache"
	"strings"
)

type Manager struct {
	executors map[int]Executor
}

type InProcExecutor struct {
	writeSetCache   wsetcache.WriteSetCache
	storageSnapshot storage.Storage
	contracts       map[string]contract.Contract
}

func (ipe *InProcExecutor) Execute(context *transaction.Context) {
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

func (ae *Manager) Execute(context *transaction.Context) {
	ae.executors[context.TX.GetExecutorType()].Execute(context)
}
