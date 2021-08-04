package wsetcache

import (
	"github.com/GrapeBaBa/brynhild/pkg/transaction"
)

type WriteSetCache interface {
	PutState(key string, write transaction.KVWrite)
	GetState(key string) transaction.KVWrite
}

func NewWriteSetCache(kind string) WriteSetCache {
	switch kind {
	case "HashMap":
		return &HashMapWriteSetCache{}
	default:
		return &HashMapWriteSetCache{}
	}
}

