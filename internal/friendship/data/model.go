package data

import "time"

type Model struct {
	ID          uint64     `db:"id"`
	RequesterID string     `db:"requester_id"`
	AddresseeID string     `db:"addressee_id"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	RespondedAt *time.Time `db:"responded_at"`
}

func (Model) TableName() string {
	return "friendships"
}
