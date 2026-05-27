package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *SeasonScore) error
	Update(ctx context.Context, entity *SeasonScore) error
	FindByUserAndSeason(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) (*SeasonScore, error)
	FindByUserIDsAndSeason(ctx context.Context, userIDs []valueobject.UserID, seasonID valueobject.SeasonID) ([]SeasonScore, error)
	ResetBySeasonID(ctx context.Context, seasonID valueobject.SeasonID) error
}
