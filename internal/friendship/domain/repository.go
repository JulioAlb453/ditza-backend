package domain

import (
	"context"

	valueobject "ditza/internal/shared/domain/value-object"
)

type Repository interface {
	Create(ctx context.Context, entity *Friendship) error
	Update(ctx context.Context, entity *Friendship) error
	FindByID(ctx context.Context, id valueobject.FriendshipID) (*Friendship, error)
	FindBetweenUsers(ctx context.Context, userA, userB valueobject.UserID) (*Friendship, error)
	FindAcceptedByUserID(ctx context.Context, userID valueobject.UserID) ([]Friendship, error)
	FindPendingByAddresseeID(ctx context.Context, addresseeID valueobject.UserID) ([]Friendship, error)
}
