package domain

import valueobject "ditza/internal/shared/domain/value-object"

type FriendEntry struct {
	UserID        valueobject.UserID
	Name          string
	SeasonPoints  int
	Rank          int
	IsCurrentUser bool
}
