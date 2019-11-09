package propagatedstorage_test

import (
	"github.com/Tanax/propagatedstorage"
	"github.com/stretchr/testify/mock"
)

type TestItem struct {
	mock.Mock
	ID              string
	Version         int
	AnotherProperty string
}

func (ti *TestItem) GetCurrentVersion() int {
	return ti.Version
}

func (ti *TestItem) GetID() string {
	return ti.ID
}

func (ti *TestItem) PopulateFromItem(item propagatedstorage.Item) error {
	args := ti.Called(item)
	return args.Error(0)
}

func CopyItemProperties(copyFrom *TestItem, copyTo *TestItem) {
	copyTo.ID = copyFrom.ID
	copyTo.Version = copyFrom.Version
	copyTo.AnotherProperty = copyFrom.AnotherProperty
}
