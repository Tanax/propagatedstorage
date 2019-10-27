package propagatedstorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Service represents how our service should look.
type Service interface {
	Get(ctx context.Context, item Item) error
	Save(ctx context.Context, item Item) error
}

type service struct {
	datastore       Datastore
	service         Service
	requiredVersion int
	itemType        Type
}

// NewService creates a new instance of a propagated storage service.
// - datastore is the datastore that stores the propagated data
// - fallbackService is the propagated storage service to fall back to if version is outdated, this is most likely a HTTP service client that asks service owning the data that is propagated
// - requiredVersion is the version required for this propagation "contract", provides a way to resync data on the fly if they ever get out of sync
// - itemType is the type of the propagated item
func NewService(ctx context.Context, datastore Datastore, itemType Type, fallbackService Service, requiredVersion int) Service {
	return &service{
		datastore:       datastore,
		service:         fallbackService,
		requiredVersion: requiredVersion,
		itemType:        itemType,
	}
}

// Get retrieves propagated data based on the (propagated) item passed in. If a required version is configured and there's a mismatch, it will try to re-fetch the item from the owning
// service if a fallback service is configured.
func (s *service) Get(ctx context.Context, item Item) error {
	model := NewModel(item.GetID(), s.itemType)

	if err := s.get(ctx, model); err != nil {
		if !errors.Is(err, ErrVersionOutdated) || s.service == nil {
			return fmt.Errorf("could not get item with id '%s': %w", item.GetID(), err)
		}

		if err := s.service.Get(ctx, item); err != nil {
			return fmt.Errorf("could not get item with id '%s': %w", item.GetID(), ErrFetchItemFromService.Wrap(err))
		}

		if err := s.Save(ctx, item); err != nil {
			return fmt.Errorf("could not get item with id '%s': %w", item.GetID(), err)
		}

		return nil
	}

	marshalled, err := json.Marshal(model.Item)
	if err != nil {
		return fmt.Errorf("could not get item with id '%s': %w", item.GetID(), ErrParsingItem.Wrap(err))
	}

	return item.PopulateFromItem(marshalled)
}

func (s *service) get(ctx context.Context, model *Model) error {
	if err := s.datastore.Get(ctx, model); err != nil {
		return fmt.Errorf("could not get propagated storage model: %w", ErrDatastoreFailed.Wrap(err))
	}

	if s.requiredVersion > 0 && s.requiredVersion > model.Version {
		return fmt.Errorf("version mismatch (required %d; current %d): %w", s.requiredVersion, model.Version, ErrVersionOutdated)
	}

	return nil
}

// Save stores propagated data based on the (propagated) item passed in. If the item's current version is higher than 0, we will assume it's the most current and update the version.
func (s *service) Save(ctx context.Context, item Item) error {
	model := NewModel(item.GetID(), s.itemType)
	model.Item = item

	if item.GetCurrentVersion() > 0 {
		model.Version = item.GetCurrentVersion()
	}

	if err := s.datastore.Save(ctx, model); err != nil {
		return ErrDatastoreFailed.Wrap(err)
	}

	return nil
}
