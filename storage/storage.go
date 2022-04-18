package storage

import (
	"context"
)

type Record struct {
	Key     string
	Value   string
	Version int64
}

type RecordIterator interface {
	Next() (bool, error)
	Value() Record
}

type Storage interface {
	GetAllRecords(ctx context.Context) (RecordIterator, error)
	Put(ctx context.Context, record Record) error
	Close() error
}
