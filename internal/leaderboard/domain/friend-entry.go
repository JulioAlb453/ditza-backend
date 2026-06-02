package domain

import valueobject "ditza/internal/shared/domain/value-object"

type FriendEntry struct {
	UserID        valueobject.UserID
	Alias         string
	SeasonPoints  int
	Rank          int
	IsCurrentUser bool
}
