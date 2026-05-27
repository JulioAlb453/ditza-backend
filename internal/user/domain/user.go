package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type User struct {
	ID         valueobject.UserID
	Name       string
	Email      string
	Password   string
	Timezone   string
	Coins      int
	FriendCode string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func New(name, email, passwordHash, timezone, friendCode string) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	timezone = strings.TrimSpace(timezone)

	if name == "" {
		return nil, domainerror.New("INVALID_USER_NAME", "el nombre es obligatorio", domainerror.ErrInvalidInput)
	}
	if utf8.RuneCountInString(name) > valueobject.MaxUserNameLength {
		return nil, domainerror.New("INVALID_USER_NAME", "el nombre excede la longitud máxima", domainerror.ErrInvalidInput)
	}
	if email == "" {
		return nil, domainerror.New("INVALID_EMAIL", "el correo es obligatorio", domainerror.ErrInvalidInput)
	}
	if timezone == "" {
		return nil, domainerror.New("INVALID_TIMEZONE", "la zona horaria es obligatoria", domainerror.ErrInvalidInput)
	}
	if friendCode == "" {
		return nil, domainerror.New("INVALID_FRIEND_CODE", "el código de amigo es obligatorio", domainerror.ErrInvalidInput)
	}

	now := time.Now().UTC()
	return &User{
		Name:       name,
		Email:      email,
		Password:   passwordHash,
		Timezone:   timezone,
		FriendCode: friendCode,
		Coins:      0,
		CreatedAt:  now,
		UpdatedAt:  now,
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
