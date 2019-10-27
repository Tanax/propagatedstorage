package propagatedstorage

import (
	"context"
)

// Datastore represents how our data store should look
type Datastore interface {
	// Get retrieves a propagated storage model that contains the propagated item. Populates
	// the model passed in by pointer to this method and returns an error if something goes wrong.
	Get(ctx context.Context, model *Model) error
	// Save stores a propagated storage model that contains the propagated item in the datastore. Returns
	// an error if it fails.
	Save(ctx context.Context, model *Model) error
}
