package data

import "time"

type Model struct {
	ID           uint64    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Timezone     string    `db:"timezone"`
	Coins        int       `db:"coins"`
	FriendCode   string    `db:"friend_code"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (Model) TableName() string {
	return "users"
}
