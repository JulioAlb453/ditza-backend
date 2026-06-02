package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type User struct {
	ID        valueobject.UserID
	Alias     string
	Email     string
	Password  string
	Coins     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(id valueobject.UserID, alias, email, hashedPassword string) (*User, error) {
	alias = strings.TrimSpace(alias)
	email = strings.TrimSpace(strings.ToLower(email))

	if id.IsZero() {
		return nil, domainerror.New("INVALID_USER_ID", "el id de usuario es obligatorio", domainerror.ErrInvalidInput)
	}
	if alias == "" {
		return nil, domainerror.New("INVALID_USER_ALIAS", "el alias es obligatorio", domainerror.ErrInvalidInput)
	}
	if utf8.RuneCountInString(alias) > valueobject.MaxUserAliasLength {
		return nil, domainerror.New("INVALID_USER_ALIAS", "el alias excede la longitud máxima", domainerror.ErrInvalidInput)
	}
	if email == "" {
		return nil, domainerror.New("INVALID_EMAIL", "el correo es obligatorio", domainerror.ErrInvalidInput)
	}
	if hashedPassword == "" {
		return nil, domainerror.New("INVALID_PASSWORD", "la contraseña es obligatoria", domainerror.ErrInvalidInput)
	}

	now := time.Now().UTC()
	return &User{
		ID:        id,
		Alias:     alias,
		Email:     email,
		Password:  hashedPassword,
		Coins:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) CanSpendCoins(amount int) bool {
	return amount > 0 && u.Coins >= amount
}

func (u *User) AddCoins(amount int) error {
	if amount < 0 {
		return domainerror.New("INVALID_COINS", "no se pueden agregar monedas negativas", domainerror.ErrInvalidInput)
	}
	u.Coins += amount
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) SpendCoins(amount int) error {
	if amount <= 0 {
		return domainerror.New("INVALID_COINS", "el monto a gastar debe ser mayor a cero", domainerror.ErrInvalidInput)
	}
	if !u.CanSpendCoins(amount) {
		return domainerror.New("INSUFFICIENT_COINS", "monedas insuficientes", domainerror.ErrInsufficientCoins)
	}
	u.Coins -= amount
	u.UpdatedAt = time.Now().UTC()
	return nil
}
