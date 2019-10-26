package propagatedstorage

import (
	"context"
)

// Datastore represents how our data store should look
type Datastore interface {
	Get(ctx context.Context, model *Model) error
	Save(ctx context.Context, model *Model) error
}
