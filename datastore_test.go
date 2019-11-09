package propagatedstorage_test

import (
	"context"

	"github.com/Tanax/propagatedstorage"
	"github.com/stretchr/testify/mock"
)

type TestDatastore struct {
	mock.Mock
}

func (ds *TestDatastore) Get(ctx context.Context, model *propagatedstorage.Model) error {
	args := ds.Called(ctx, model)
	if len(args) > 1 {
		if responseModel, ok := args.Get(1).(*propagatedstorage.Model); ok {
			*model = *responseModel
		}
	}
	return args.Error(0)
}

func (ds *TestDatastore) Save(ctx context.Context, model *propagatedstorage.Model) error {
	args := ds.Called(ctx, model)
	return args.Error(0)
}
