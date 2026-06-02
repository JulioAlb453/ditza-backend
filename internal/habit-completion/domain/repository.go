package domain

import (
	"context"
	"time"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *HabitCompletion) error
	ExistsForHabitOnDate(ctx context.Context, habitID valueobject.HabitID, date time.Time) (bool, error)
	CountByUserOnDate(ctx context.Context, userID valueobject.UserID, date time.Time) (int, error)
	FindByUserIDAndDateRange(ctx context.Context, userID valueobject.UserID, from, to time.Time) ([]HabitCompletion, error)
}
