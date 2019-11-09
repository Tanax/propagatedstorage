package documentstore

import (
	"time"

	"github.com/Tanax/propagatedstorage"
)

// Entity defines how our documents store entity looks like.
type Entity struct {
	ID      string
	Type    propagatedstorage.Type
	Version int
	Item    propagatedstorage.Item

	Created  time.Time
	Modified time.Time
}

func (e *Entity) populateModel(model *propagatedstorage.Model) error {
	model.ID = e.ID
	model.Type = e.Type
	model.Item = e.Item
	model.Version = e.Version
	model.Created = e.Created
	model.Modified = e.Modified

	return nil
}

// NewFromModel creates a new Entity based on a Model.
func NewFromModel(model *propagatedstorage.Model) (*Entity, error) {
	e := new(Entity)

	e.ID = model.ID
	e.Type = model.Type
	e.Item = model.Item
	e.Version = model.Version
	e.Created = model.Created
	e.Modified = model.Modified

	return e, nil
}
