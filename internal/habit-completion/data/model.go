package data

import "time"

type Model struct {
	ID                  uint64    `db:"id"`
	HabitID             uint64    `db:"habit_id"`
	UserID              string    `db:"user_id"`
	CompletedAt         time.Time `db:"completed_at"`
	Note                *string   `db:"note"`
	Emoji               *string   `db:"emoji"`
	CoinsAwarded        int       `db:"coins_awarded"`
	SeasonPointsAwarded int       `db:"season_points_awarded"`
	CreatedAt           time.Time `db:"created_at"`
}

func (Model) TableName() string {
	return "habit_completions"
}
