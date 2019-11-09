package propagatedstorage

// Item describes how a propagated data item should look
type Item interface {
	GetCurrentVersion() int
	GetID() string
	PopulateFromItem(item Item) error
}
