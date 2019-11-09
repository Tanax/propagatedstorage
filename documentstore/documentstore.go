package documentstore

import (
	"context"

	"github.com/Tanax/propagatedstorage"
)

type documentstore struct {
	coll Collection
}

// New returns a new propagated storage datastore of type document store
func New(coll Collection) propagatedstorage.Datastore {
	return &documentstore{
		coll: coll,
	}
}

func (ds *documentstore) Get(ctx context.Context, model *propagatedstorage.Model) error {
	entity, err := NewFromModel(model)
	if err != nil {
		return err
	}

	if err := ds.coll.Get(ctx, entity); err != nil {
		return err
	}

	return entity.populateModel(model)
}

func (ds *documentstore) Save(ctx context.Context, model *propagatedstorage.Model) error {
	entity, err := NewFromModel(model)
	if err != nil {
		return err
	}

	if err := ds.coll.Put(ctx, entity); err != nil {
		return err
	}

	return entity.populateModel(model)
}
