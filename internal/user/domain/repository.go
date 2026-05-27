package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *User) error
	Update(ctx context.Context, entity *User) error
	FindByID(ctx context.Context, id valueobject.UserID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
