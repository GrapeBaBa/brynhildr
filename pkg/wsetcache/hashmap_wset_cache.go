package wsetcache

import (
	"github.com/GrapeBaBa/brynhildr/pkg/transaction"
	"sync"
)

type HashMapWriteSetCache struct {
	cache sync.Map
}

func (hmwsc *HashMapWriteSetCache) PutState(key string, value transaction.KVWrite) {
	hmwsc.cache.Store(key, value)
}

func (hmwsc *HashMapWriteSetCache) GetState(key string) transaction.KVWrite {
	value, ok := hmwsc.cache.Load(key)
	if !ok {
		return transaction.KVWrite{}
	} else {
		return value.(transaction.KVWrite)
	}

}
