package store

import (
	"mxshs/redis-go/app/types"
	"sync"
	"time"
)

type kvValue struct {
	Value  *types.Data
	Expiry int64
}

type KVStore struct {
    //store map[string]*kvValue
    store *hashmap[kvValue]
    sync.RWMutex
}

func NewKVStore() *KVStore {
	return &KVStore{
        store: NewHashmap[kvValue](10000),
        //store: make(map[string]*kvValue),
    }
}

func (kv *KVStore) Set(key string, value *types.Data, expiry int64) error {
	kv.store.Set(key, kvValue{
		Value:  value,
		Expiry: expiry,
	})

	return nil
}

func (kv *KVStore) Get(key string) (*types.Data, error) {
	value, ok := kv.store.Get(key)
	if !ok {
		return nil, types.NotFound
	}

	if value.Expiry > 0 && time.Now().UnixMilli() > value.Expiry {
		// TODO: implement some sort of a queue or something to evict before an expired key is first requested (on 10k connections writing to random keys it grows too fast)
        //delete(kv.store, key)
        kv.store.Delete(key)
		return nil, types.Expired
	}

	return value.Value, nil
}

