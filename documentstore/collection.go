package documentstore

import (
	"context"

	"gocloud.dev/docstore"
)

// Collection defines the collection interface that we interact with within this documentstore.
type Collection interface {
	Get(ctx context.Context, document interface{}, fps ...docstore.FieldPath) error
	Put(ctx context.Context, document interface{}) error
}
