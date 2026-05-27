package data

import "time"

type Model struct {
	ID          uint64    `db:"id"`
	UserID      uint64    `db:"user_id"`
	Type        string    `db:"type"`
	CoinsDelta  int       `db:"coins_delta"`
	SeasonDelta int       `db:"season_delta"`
	ReferenceID *uint64   `db:"reference_id"`
	CreatedAt   time.Time `db:"created_at"`
}

func (Model) TableName() string {
	return "point_transactions"
}
