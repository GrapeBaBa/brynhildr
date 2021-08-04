package committer

import (
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
)

type Committer interface {
	Commit(ctxs []*transaction.Context)
}
