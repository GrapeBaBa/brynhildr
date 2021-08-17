package executor

import (
	"github.com/GrapeBaBa/brynhildr/pkg/contract"
	"github.com/GrapeBaBa/brynhildr/pkg/storage"
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"strings"
)

type InProcTransactionExecutor struct {
	storage   storage.Storage
	contracts map[string]contract.Contract
}

func NewInProcTransactionExecutor(storage storage.Storage, contracts map[string]contract.Contract) *InProcTransactionExecutor {
	return &InProcTransactionExecutor{
		storage:   storage,
		contracts: contracts,
	}
}

func (ipe *InProcTransactionExecutor) Execute(context *transaction.Context) {
	tx := context.Transaction
	//TODO:Prevent contract panic
	ipccs := contract.NewInProcContractCallStub(context, ipe.storage)
	ctc := contract.NewInProcContractTransactionContext(ipccs)
	// no use reflection for performance consideration
	var execRes transaction.ExecResult
	if strings.EqualFold(tx.GetMethod(), "init") {
		execRes = ipe.contracts[tx.GetContractID()].Init(ctc)
	} else {
		execRes = ipe.contracts[tx.GetContractID()].Invoke(ctc)
	}
	context.Result.ExecResult = execRes
}


