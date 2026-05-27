package data

import "time"

type Model struct {
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	Slot       string    `db:"slot"`
	PriceCoins int       `db:"price_coins"`
	Rarity     string    `db:"rarity"`
	AssetKey   string    `db:"asset_key"`
	IsActive   bool      `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
}

func (Model) TableName() string {
	return "cosmetics"
}
