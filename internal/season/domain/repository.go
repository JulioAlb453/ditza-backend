package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *Season) error
	Update(ctx context.Context, entity *Season) error
	FindByID(ctx context.Context, id valueobject.SeasonID) (*Season, error)
	FindActive(ctx context.Context) (*Season, error)
	DeactivateAll(ctx context.Context) error
}
