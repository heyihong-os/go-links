package main

import (
	"context"
	"testing"

	"github.com/heyihong-os/go-links/storage"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	ctx   context.Context
	store *LinkStore
}

func (suite *StoreTestSuite) SetupTest() {
	suite.ctx = context.Background()
	store, err := NewLinkStore(suite.ctx, storage.NewInMemStorage())
	if err != nil {
		suite.Fail("Fail to create link store: %v", err)
	}
	suite.store = store
}

func (suite *StoreTestSuite) TestBasic() {
	suite.Nil(suite.store.PutLink(suite.ctx, "test", "https://test"))

	link, err := suite.store.GetLink(suite.ctx, "test")
	suite.Nil(err)
	suite.Equal("https://test", *link)

	suite.Nil(suite.store.PutLink(suite.ctx, "test", "https://test1"))

	link, err = suite.store.GetLink(suite.ctx, "test")
	suite.Nil(err)
	suite.Equal("https://test1", *link)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
