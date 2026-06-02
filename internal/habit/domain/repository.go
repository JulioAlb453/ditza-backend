package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *Habit) error
	Update(ctx context.Context, entity *Habit) error
	FindByID(ctx context.Context, id valueobject.HabitID) (*Habit, error)
	FindActiveByUserID(ctx context.Context, userID valueobject.UserID) ([]Habit, error)
	CountActiveByUserID(ctx context.Context, userID valueobject.UserID) (int, error)
}
