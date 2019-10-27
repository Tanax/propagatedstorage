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
func NewModel(id string, itemType Type) *Model {
	model := new(Model)
	model.ID = id
	model.Type = itemType

	return model
}

// Item describes how a propagated data item should look
type Item interface {
	GetCurrentVersion() int
	GetID() string
	GetType() Type
	PopulateFromItem(item []byte) error
}

// Type declares how a propagated item type looks
type Type string
