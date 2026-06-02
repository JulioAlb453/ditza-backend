package data

import "time"

type Model struct {
	ID        string    `db:"id"`
	Alias     string    `db:"alias"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Coins     int       `db:"coins"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (Model) TableName() string {
	return "users"
}
