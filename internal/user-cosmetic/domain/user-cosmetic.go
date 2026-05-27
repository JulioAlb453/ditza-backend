package domain

import (
	"time"

	valueobject "ditza/internal/shared/domain/value-object"
)

type UserCosmetic struct {
	UserID      valueobject.UserID
	CosmeticID  valueobject.CosmeticID
	PurchasedAt time.Time
}

func New(userID valueobject.UserID, cosmeticID valueobject.CosmeticID) *UserCosmetic {
	return &UserCosmetic{
		UserID:      userID,
		CosmeticID:  cosmeticID,
		PurchasedAt: time.Now().UTC(),
	}
}
