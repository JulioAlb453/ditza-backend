package domain

import "fmt"

type Type string

const (
	TypeHabit       Type = "habit"
	TypeStreakBonus Type = "streak_bonus"
	TypeNoteBonus   Type = "note_bonus"
	TypePurchase    Type = "purchase"
	TypeSeasonReset Type = "season_reset"
)

func ParseType(raw string) (Type, error) {
	txType := Type(raw)
	switch txType {
	case TypeHabit, TypeStreakBonus, TypeNoteBonus, TypePurchase, TypeSeasonReset:
		return txType, nil
	default:
		return "", fmt.Errorf("tipo de transacción de puntos inválido: %s", raw)
	}
}

func (t Type) String() string {
	return string(t)
}

func (t Type) IsValid() bool {
	_, err := ParseType(string(t))
	return err == nil
}
