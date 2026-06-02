package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *UserCosmetic) error
	Exists(ctx context.Context, userID valueobject.UserID, cosmeticID valueobject.CosmeticID) (bool, error)
	FindByUserID(ctx context.Context, userID valueobject.UserID) ([]UserCosmetic, error)
}
