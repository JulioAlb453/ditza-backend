package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *PointTransaction) error
	FindByUserID(ctx context.Context, userID valueobject.UserID, limit int) ([]PointTransaction, error)
}
