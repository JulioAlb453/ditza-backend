package data

import "time"

type Model struct {
	ID        uint64    `db:"id"`
	StartsAt  time.Time `db:"starts_at"`
	EndsAt    time.Time `db:"ends_at"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}

func (Model) TableName() string {
	return "seasons"
}
