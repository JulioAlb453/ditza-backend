package data

import "time"

type Model struct {
	ID                uint64     `db:"id"`
	UserID            string     `db:"user_id"`
	Title             string     `db:"title"`
	Description       string     `db:"description"`
	Emoji             string     `db:"emoji"`
	Category          string     `db:"category"`
	Color             string     `db:"color"`
	Frequency         string     `db:"frequency"`
	TargetCount       int        `db:"target_count"`
	TargetUnit        string     `db:"target_unit"`
	Difficulty        string     `db:"difficulty"`
	ReminderTime      *string    `db:"reminder_time"`
	IsActive          bool       `db:"is_active"`
	CurrentStreak     int        `db:"current_streak"`
	BestStreak        int        `db:"best_streak"`
	LastCompletedDate *time.Time `db:"last_completed_date"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

func (Model) TableName() string {
	return "habits"
}
