package documentstore_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tanax/propagatedstorage"
	"github.com/Tanax/propagatedstorage/documentstore"
	"github.com/stretchr/testify/assert"
)

func TestGet_Success(t *testing.T) {
	// Setup
	ctx := context.TODO()

	// Mock
	var (
		model          = &propagatedstorage.Model{ID: "ThisIsMyID"}
		inputEntity, _ = documentstore.NewFromModel(model)
		responseEntity = &documentstore.Entity{ID: "AnotherID", Type: "MyType", Version: 3, Created: time.Now(), Modified: time.Now()}
		collection     = &TestCollection{}
	)

	// Expect
	collection.On("Get", ctx, inputEntity).Return(nil, responseEntity)

	// Apply
	docstore := documentstore.New(collection)
	err := docstore.Get(ctx, model)

	// Assert
	collection.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, responseEntity.ID, model.ID)
	assert.Equal(t, responseEntity.Type, model.Type)
	assert.Equal(t, responseEntity.Version, model.Version)
	assert.Equal(t, responseEntity.Created, model.Created)
	assert.Equal(t, responseEntity.Modified, model.Modified)
}

func TestGet_CollectionError(t *testing.T) {
	// Setup
	ctx := context.TODO()

	// Mock
	var (
		model          = &propagatedstorage.Model{ID: "ThisIsMyID"}
		inputEntity, _ = documentstore.NewFromModel(model)
		collection     = &TestCollection{}
	)

	// Expect
	collection.On("Get", ctx, inputEntity).Return(errors.New("error"))

	// Apply
	docstore := documentstore.New(collection)
	err := docstore.Get(ctx, model)

	// Assert
	collection.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "error")
}

func TestSave_Success(t *testing.T) {
	// Setup
	ctx := context.TODO()

	// Mock
	var (
		model          = &propagatedstorage.Model{ID: "ThisIsMyID"}
		inputEntity, _ = documentstore.NewFromModel(model)
		responseEntity = &documentstore.Entity{ID: "AnotherID"}
		collection     = &TestCollection{}
	)

	// Expect
	collection.On("Put", ctx, inputEntity).Return(nil, responseEntity)

	// Apply
	docstore := documentstore.New(collection)
	err := docstore.Save(ctx, model)

	// Assert
	collection.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, responseEntity.ID, model.ID)
}

func TestSave_CollectionError(t *testing.T) {
	// Setup
	ctx := context.TODO()

	// Mock
	var (
		model          = &propagatedstorage.Model{ID: "ThisIsMyID"}
		inputEntity, _ = documentstore.NewFromModel(model)
		collection     = &TestCollection{}
	)

	// Expect
	collection.On("Put", ctx, inputEntity).Return(errors.New("error"))

	// Apply
	docstore := documentstore.New(collection)
	err := docstore.Save(ctx, model)

	// Assert
	collection.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "error")
}
