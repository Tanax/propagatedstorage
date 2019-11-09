package propagatedstorage_test

import "github.com/Tanax/propagatedstorage"

func MockModelWithItem(item *TestItem, itemType propagatedstorage.Type) *propagatedstorage.Model {
	model := MockModel(item, itemType)
	model.Item = item
	return model
}

func MockModel(item *TestItem, itemType propagatedstorage.Type) *propagatedstorage.Model {
	return propagatedstorage.NewModel(item.ID, itemType, item.Version)
}
