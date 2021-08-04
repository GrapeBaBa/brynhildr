package executor

import (
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
)

type Executor interface {
	Execute(context *transaction.Context)
}
