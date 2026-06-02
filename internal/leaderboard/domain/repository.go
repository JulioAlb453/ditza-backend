package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	GetFriendRanking(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) ([]FriendEntry, error)
}
