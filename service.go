package propagatedstorage

import (
	"context"
	"fmt"
)

// Service represents how our service should look.
type Service interface {
	Get(ctx context.Context, item Item) error
	Save(ctx context.Context, item Item) error
}

type service struct {
	datastore       Datastore
	requiredVersion int
	itemType        Type
	fallbackService Service
}

// NewService creates a new instance of a propagated storage service.
// - datastore is the datastore that stores the propagated data
// - fallbackService is the propagated storage service to fall back to if version is outdated, this is most likely a HTTP service client that asks service owning the data that is propagated
// - requiredVersion is the version required for this propagation "contract", provides a way to resync data on the fly if they ever get out of sync
// - itemType is the type of the propagated item
func NewService(datastore Datastore, itemType Type, requiredVersion int, fallbackService Service) Service {
	return &service{
		datastore:       datastore,
		requiredVersion: requiredVersion,
		itemType:        itemType,
		fallbackService: fallbackService,
	}
}

// Get retrieves propagated data based on the (propagated) item passed in. If a required version is configured, it will check that against what was stored in our propagated storage
// and return an error if it's below the required version.
func (s *service) Get(ctx context.Context, item Item) error {
	model := NewModel(item.GetID(), s.itemType, item.GetCurrentVersion())

	if err := s.datastore.Get(ctx, model); err != nil {
		return fmt.Errorf("could not get propagated storage model: %w", ErrDatastoreFailed.Wrap(err))
	}

	if err := s.validateVersion(model.Version); err != nil || model.Item == nil {
		if s.fallbackService == nil {
			return fmt.Errorf("could not get propagated item from fallback service: %w", ErrMissingFallbackService.Wrap(err))
		}

		if err := s.fallbackService.Get(ctx, model.Item); err != nil {
			return fmt.Errorf("could not get propagated item from fallback service: %w", ErrServiceFailed.Wrap(err))
		}

		if err := s.Save(ctx, model.Item); err != nil {
			return fmt.Errorf("could not save propagated item from fallback service: %w", err)
		}
	}

	return item.PopulateFromItem(model.Item)
}

func (s *service) validateVersion(version int) error {
	if s.requiredVersion > 0 && s.requiredVersion > version {
		return fmt.Errorf("version mismatch (required %d; current %d): %w", s.requiredVersion, version, ErrVersionOutdated)
	}
	return nil
}

// Save stores propagated data based on the (propagated) item passed in. If the item's current version is higher than 0, we will assume it's the most current and update the version.
func (s *service) Save(ctx context.Context, item Item) error {
	model := NewModel(item.GetID(), s.itemType, item.GetCurrentVersion())
	model.Item = item

	if err := s.datastore.Save(ctx, model); err != nil {
		return ErrDatastoreFailed.Wrap(err)
	}

	return nil
}
