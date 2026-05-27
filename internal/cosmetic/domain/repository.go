package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *Cosmetic) error
	Update(ctx context.Context, entity *Cosmetic) error
	FindByID(ctx context.Context, id valueobject.CosmeticID) (*Cosmetic, error)
	FindAllActive(ctx context.Context) ([]Cosmetic, error)
}
