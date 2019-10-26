package dynamodb

import (
	"context"
	"fmt"
	"sync"

	"github.com/Tanax/propagatedstorage"
	"github.com/Tanax/propagatedstorage/documentstore"
	"github.com/aws/aws-sdk-go/aws/session"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"gocloud.dev/docstore/awsdynamodb"
)

// InitiateAsync initializes a propagated storage datastore with a dynamo db driver asynchronously.
func InitiateAsync(ctx context.Context, sess *session.Session, datastore *propagatedstorage.Datastore, tableName string, errors []error, wg *sync.WaitGroup) {
	docstore, err := InitiateSync(ctx, sess, tableName)
	if err != nil {
		errors = append(errors, err)
	} else {
		*datastore = docstore
	}

	wg.Done()
}

// InitiateSync initializes a propagated storage datastore with a dynamo db driver synchronously.
func InitiateSync(ctx context.Context, sess *session.Session, tableName string) (propagatedstorage.Datastore, error) {
	if sess == nil {
		return nil, fmt.Errorf("failed to open collection propagated storage: %w", ErrMissingSession)
	}

	if tableName == "" {
		tableName = "propagatedstorage"
	}

	driver, err := awsdynamodb.OpenCollection(ddb.New(sess), tableName, "Type", "ID", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open collection propagated storage: %w", ErrFailedToInitiate.Wrap(err))
	}

	return documentstore.New(driver), nil
}
