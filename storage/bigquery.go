package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type BigQueryRecordIterator struct {
	record Record
	bi     *bigquery.RowIterator
}

func NewBigQueryRecordIterator(bi *bigquery.RowIterator) *BigQueryRecordIterator {
	return &BigQueryRecordIterator{
		bi: bi,
	}
}

func (ri *BigQueryRecordIterator) Next() (bool, error) {
	err := ri.bi.Next(&ri.record)
	if err != nil {
		if err == iterator.Done {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ri *BigQueryRecordIterator) Value() Record {
	return ri.record
}

var _ RecordIterator = (*BigQueryRecordIterator)(nil)

type BigQueryStorage struct {
	bqClient    *bigquery.Client
	datasetName string
	tableName   string
}

var _ Storage = (*BigQueryStorage)(nil)

func NewBigQueryStorage(
	bqClient *bigquery.Client,
	datasetName string,
	tableName string,
) *BigQueryStorage {
	return &BigQueryStorage{
		bqClient:    bqClient,
		datasetName: datasetName,
		tableName:   tableName,
	}
}

func (b *BigQueryStorage) GetAllRecords(ctx context.Context) (RecordIterator, error) {
	q := b.bqClient.Query(fmt.Sprintf("SELECT * FROM %s.%s", b.datasetName, b.tableName))
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	return NewBigQueryRecordIterator(it), nil
}

func (b *BigQueryStorage) Put(ctx context.Context, row Record) error {
	ins := b.bqClient.Dataset(b.datasetName).Table(b.tableName).Inserter()

	return ins.Put(ctx, []*Record{&row})
}

func (b *BigQueryStorage) Close() error {
	return b.bqClient.Close()
}
