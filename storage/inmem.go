package storage

import (
	"context"
	"sync"
)

type InMemIterator struct {
	idx     int
	records []Record
}

func NewInMemIterator(records []Record) *InMemIterator {
	return &InMemIterator{
		idx:     -1,
		records: records,
	}
}

func (imi *InMemIterator) Next() (bool, error) {
	if imi.idx+1 >= len(imi.records) {
		return false, nil
	}
	imi.idx++
	return true, nil
}

func (imi *InMemIterator) Value() Record {
	return imi.records[imi.idx]
}

var _ RecordIterator = (*InMemIterator)(nil)

type InMemStorage struct {
	mu      sync.Mutex
	records []Record
}

var _ Storage = (*InMemStorage)(nil)

func NewInMemStorage() *InMemStorage {
	return &InMemStorage{}
}

func (ims *InMemStorage) GetAllRecords(ctx context.Context) (RecordIterator, error) {
	ims.mu.Lock()
	defer ims.mu.Unlock()

	return NewInMemIterator(ims.records), nil
}

func (ims *InMemStorage) Put(ctx context.Context, record Record) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()

	ims.records = append(ims.records, record)
	return nil
}

func (*InMemStorage) Close() error {
	return nil
}
