package propagatedstorage

import (
	"fmt"
)

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

func (e *BaseError) Unwrap() error {
	return e.err
}

func (e *BaseError) Wrap(err error) *BaseError {
	e.err = err
	return e
}

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
)
