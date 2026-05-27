package data

import "time"

type Model struct {
	UserID      string    `db:"user_id"`
	CosmeticID  uint64    `db:"cosmetic_id"`
	PurchasedAt time.Time `db:"purchased_at"`
}

func (Model) TableName() string {
	return "user_cosmetics"
}
