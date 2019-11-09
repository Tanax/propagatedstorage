package propagatedstorage

import (
	"fmt"
)

// BaseError ..
type BaseError struct {
	err error
	msg string
}

func (e *BaseError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}
	return e.msg
}

// Unwrap ..
func (e *BaseError) Unwrap() error {
	return e.err
}

// Wrap ..
func (e *BaseError) Wrap(err error) *BaseError {
	e.err = err
	return e
}

// NewError ..
func NewError(message string) *BaseError {
	err := new(BaseError)
	err.msg = message

	return err
}

var (
	// ErrVersionOutdated ..
	ErrVersionOutdated = NewError("version outdated")
	// ErrDatastoreFailed ..
	ErrDatastoreFailed = NewError("datastore failed")
	// ErrMissingFallbackService ...
	ErrMissingFallbackService = NewError("missing fallback service")
	// ErrServiceFailed ..
	ErrServiceFailed = NewError("service failed")
	// ErrMissingDatastoreSession ..
	ErrMissingDatastoreSession = NewError("session not found")
	// ErrInitiateDatastoreDriver ..
	ErrInitiateDatastoreDriver = NewError("failed to initiate driver")
	// ErrFetchItemFromService ..
	ErrFetchItemFromService = NewError("failed to fetch item from fallback service")
)
