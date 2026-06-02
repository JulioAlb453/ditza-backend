package data

import "time"

type Model struct {
	UserID    string    `db:"user_id"`
	SeasonID  uint64    `db:"season_id"`
	Points    int       `db:"points"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (Model) TableName() string {
	return "season_scores"
}
