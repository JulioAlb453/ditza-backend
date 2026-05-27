package data

import "time"

type Model struct {
	UserID               uint64    `db:"user_id"`
	Name                 string    `db:"name"`
	Level                int       `db:"level"`
	XP                   int       `db:"xp"`
	Mood                 string    `db:"mood"`
	EquippedHatID        *uint64   `db:"equipped_hat_id"`
	EquippedShirtID      *uint64   `db:"equipped_shirt_id"`
	EquippedBackgroundID *uint64   `db:"equipped_background_id"`
	EquippedAccessoryID  *uint64   `db:"equipped_accessory_id"`
	LastInteractionAt    time.Time `db:"last_interaction_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}

func (Model) TableName() string {
	return "pets"
}
