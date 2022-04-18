package main

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/heyihong-os/go-links/storage"
)

type Value struct {
	Value   string
	Version int64
}

type KeyValues interface {
	Get(key string) (Value, bool)
	Put(key string, val Value)
}

type KeyValueMap struct {
	rwLock sync.RWMutex
	kvMap  map[string]Value
}

func NewKeyValueMap() *KeyValueMap {
	return &KeyValueMap{
		kvMap: make(map[string]Value),
	}
}

func (kvm *KeyValueMap) Get(key string) (Value, bool) {
	kvm.rwLock.RLock()
	defer kvm.rwLock.RUnlock()

	v, ok := kvm.kvMap[key]
	if !ok {
		return Value{}, false
	}
	return v, true
}

func (kvm *KeyValueMap) Put(key string, val Value) {
	kvm.rwLock.Lock()
	defer kvm.rwLock.Unlock()

	v, ok := kvm.kvMap[key]
	if !ok || (v.Version < val.Version || (v.Version == val.Version && v.Value < val.Value)) {
		kvm.kvMap[key] = val
	}
}

type LinkStore struct {
	kvs     KeyValues
	storage storage.Storage
}

func NewLinkStore(ctx context.Context, storage storage.Storage) (*LinkStore, error) {
	iter, err := storage.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	kvs := NewKeyValueMap()

	for {
		hasNext, err := iter.Next()
		if err != nil {
			return nil, err
		}
		if !hasNext {
			break
		}
		row := iter.Value()
		kvs.Put(row.Key, Value{
			Value:   row.Value,
			Version: row.Version,
		})
	}

	return &LinkStore{
		kvs:     kvs,
		storage: storage,
	}, nil
}

func (ls *LinkStore) GetLink(ctx context.Context, shortLink string) (*string, error) {
	v, ok := ls.kvs.Get(toKey(shortLink))
	if !ok {
		return nil, nil
	}
	return &v.Value, nil
}

func (ls *LinkStore) PutLink(ctx context.Context, shortLink string, originalLink string) error {
	key := toKey(shortLink)
	newVal := Value{
		Value:   originalLink,
		Version: time.Now().Unix(),
	}
	if curVal, ok := ls.kvs.Get(key); ok {
		newVal.Version = maxInt64(newVal.Version, curVal.Version+1)
	}

	if err := ls.storage.Put(
		ctx,
		storage.Record{Key: key, Value: newVal.Value, Version: newVal.Version},
	); err != nil {
		return err
	}

	ls.kvs.Put(key, newVal)
	return nil
}

func maxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func toKey(shortLink string) string {
	return strings.ReplaceAll(shortLink, "/", "")
}
