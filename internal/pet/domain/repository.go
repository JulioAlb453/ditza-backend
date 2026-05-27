package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *Pet) error
	Update(ctx context.Context, entity *Pet) error
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*Pet, error)
}
