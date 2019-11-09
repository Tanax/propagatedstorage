package propagatedstorage

import (
	"time"
)

// Model represents how our propagated data model should look
type Model struct {
	ID      string
	Type    Type
	Version int
	Item    Item

	Created  time.Time
	Modified time.Time
}

// NewModel creates a new propagated data model
func NewModel(id string, itemType Type, version int) *Model {
	model := new(Model)
	model.ID = id
	model.Type = itemType
	model.Version = version

	return model
}
