package documentstore_test

import (
	"context"

	"github.com/Tanax/propagatedstorage/documentstore"
	"github.com/stretchr/testify/mock"
	"gocloud.dev/docstore"
)

type TestCollection struct {
	mock.Mock
}

func (tc *TestCollection) Get(ctx context.Context, document interface{}, fps ...docstore.FieldPath) error {
	args := tc.Called(ctx, document)
	if len(args) > 1 {
		if entity, ok := args.Get(1).(*documentstore.Entity); ok {
			if inputEntity, ok := document.(*documentstore.Entity); ok {
				*entity = *inputEntity
			}
		}
	}
	return args.Error(0)
}

func (tc *TestCollection) Put(ctx context.Context, document interface{}) error {
	args := tc.Called(ctx, document)
	if len(args) > 1 {
		if entity, ok := args.Get(1).(*documentstore.Entity); ok {
			if inputEntity, ok := document.(*documentstore.Entity); ok {
				*entity = *inputEntity
			}
		}
	}
	return args.Error(0)
}
