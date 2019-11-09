package propagatedstorage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Tanax/propagatedstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestService struct {
	mock.Mock
}

func (ts *TestService) Get(ctx context.Context, item propagatedstorage.Item) error {
	args := ts.Called(ctx, item)
	if len(args) > 1 {
		if responseItem, ok := args.Get(1).(*TestItem); ok {
			if inputItem, ok := item.(*TestItem); ok {
				CopyItemProperties(responseItem, inputItem)
			}
		}
	}
	return args.Error(0)
}

func (ts *TestService) Save(ctx context.Context, item propagatedstorage.Item) error {
	args := ts.Called(ctx, item)
	return args.Error(0)
}

func TestGet_PopulateFromItemError(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem    = &TestItem{ID: testId}
		responseItem = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 0}
		datastore    = &TestDatastore{}
	)

	// Expect
	datastore.On("Get", ctx, MockModel(inputItem, testType)).Return(nil, MockModelWithItem(responseItem, testType))
	inputItem.On("PopulateFromItem", responseItem).Return(errors.New("error"))

	// Apply
	service := propagatedstorage.NewService(datastore, TestType, 0, nil)
	err := service.Get(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "error")
}

func TestGet_DatastoreError(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem = &TestItem{ID: testId}
		datastore = &TestDatastore{}
	)

	// Expect
	datastore.On("Get", ctx, MockModel(inputItem, testType)).Return(errors.New("error"))

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 1, nil)
	err := service.Get(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, propagatedstorage.ErrDatastoreFailed))
}

func TestGet_VersionMismatch(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem    = &TestItem{ID: testId}
		responseItem = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 0}
		datastore    = &TestDatastore{}
	)

	// Expect
	datastore.On("Get", ctx, MockModel(inputItem, testType)).Return(nil, MockModelWithItem(responseItem, testType))

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 1, nil)
	err := service.Get(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, propagatedstorage.ErrVersionOutdated))
	assert.True(t, errors.Is(err, propagatedstorage.ErrMissingFallbackService))
}

func TestGet_FallbackServiceFailed(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem       = &TestItem{ID: testId}
		responseItem    = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 0}
		datastore       = &TestDatastore{}
		fallbackService = &TestService{}
	)

	// Expect
	datastore.On("Get", ctx, MockModel(inputItem, testType)).Return(nil, MockModelWithItem(responseItem, testType))
	fallbackService.On("Get", ctx, responseItem).Return(errors.New("error"))

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 1, fallbackService)
	err := service.Get(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, propagatedstorage.ErrServiceFailed))
}

func TestGet_UpdatePropagatedItemFailed(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem             = &TestItem{ID: testId}
		datastoreResponseItem = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 0}
		serviceResponseItem   = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 1}
		datastore             = &TestDatastore{}
		fallbackService       = &TestService{}
	)

	// Expect
	datastore.On("Get", ctx, MockModel(inputItem, testType)).Return(nil, MockModelWithItem(datastoreResponseItem, testType))
	fallbackService.On("Get", ctx, datastoreResponseItem).Return(nil, serviceResponseItem)
	datastore.On("Save", ctx, MockModelWithItem(serviceResponseItem, testType)).Return(errors.New("error"))

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 1, fallbackService)
	err := service.Get(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, propagatedstorage.ErrDatastoreFailed))
}

func TestGet_UpdatePropagatedItemSuccess(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem             = &TestItem{ID: testId}
		datastoreResponseItem = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 0}
		serviceResponseItem   = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 1}
		datastore             = &TestDatastore{}
		fallbackService       = &TestService{}
	)

	// Expect
	datastore.On("Get", ctx, MockModel(inputItem, testType)).Return(nil, MockModelWithItem(datastoreResponseItem, testType))
	fallbackService.On("Get", ctx, datastoreResponseItem).Return(nil, serviceResponseItem)
	datastore.On("Save", ctx, MockModelWithItem(serviceResponseItem, testType)).Return(nil)
	inputItem.On("PopulateFromItem", serviceResponseItem).Return(nil)

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 1, fallbackService)
	err := service.Get(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.Nil(t, err)
}

func TestGet_Success(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem    = &TestItem{ID: testId}
		responseItem = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 0}
		datastore    = &TestDatastore{}
	)

	// Expect
	datastore.On("Get", ctx, MockModel(inputItem, testType)).Return(nil, MockModelWithItem(responseItem, testType))
	inputItem.On("PopulateFromItem", responseItem).Return(nil)

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 0, nil)
	err := service.Get(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.Nil(t, err)
}

func TestSave_DatastoreError(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 2}
		datastore = &TestDatastore{}
	)

	// Expect
	datastore.On("Save", ctx, MockModelWithItem(inputItem, testType)).Return(errors.New("error"))

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 0, nil)
	err := service.Save(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, propagatedstorage.ErrDatastoreFailed))
}

func TestSave_Success(t *testing.T) {
	// Setup
	var (
		testId   = "ThisIsMyID"
		testType = TestType
		ctx      = context.TODO()
	)

	// Mock
	var (
		inputItem = &TestItem{ID: testId, AnotherProperty: "Heyhey", Version: 2}
		datastore = &TestDatastore{}
	)

	// Expect
	datastore.On("Save", ctx, MockModelWithItem(inputItem, testType)).Return(nil)

	// Apply
	service := propagatedstorage.NewService(datastore, testType, 0, nil)
	err := service.Save(ctx, inputItem)

	// Assert
	inputItem.AssertExpectations(t)
	datastore.AssertExpectations(t)
	assert.Nil(t, err)
}
