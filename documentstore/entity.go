package documentstore

import (
	"time"

	"github.com/Tanax/propagatedstorage"
)

type entity struct {
	ID      string
	Type    propagatedstorage.Type
	Version int
	Item    propagatedstorage.Item

	Created  time.Time
	Modified time.Time
}

func (e *entity) populateModel(model *propagatedstorage.Model) error {
	model.ID = e.ID
	model.Type = e.Type
	model.Item = e.Item
	model.Version = e.Version
	model.Created = e.Created
	model.Modified = e.Modified

	return nil
}

func newEntityFromModel(model *propagatedstorage.Model) (*entity, error) {
	e := new(entity)

	e.ID = model.ID
	e.Type = model.Type
	e.Item = model.Item
	e.Version = model.Version
	e.Created = model.Created
	e.Modified = model.Modified

	return e, nil
}
