package data

import (
	"ditza/internal/user/domain"
	valueobject "ditza/internal/shared/domain/value-object"
)
func ToDomain(m Model) domain.User {
	return domain.User{
		ID:         valueobject.UserID(m.ID),
		Name:       m.Name,
		Email:      m.Email,
		Password:   m.PasswordHash,
		Timezone:   m.Timezone,
		Coins:      m.Coins,
		FriendCode: m.FriendCode,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func ToModel(e domain.User) Model {
	return Model{
		ID:           uint64(e.ID),
		Name:         e.Name,
		Email:        e.Email,
		PasswordHash: e.Password,
		Timezone:     e.Timezone,
		Coins:        e.Coins,
		FriendCode:   e.FriendCode,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func ToDomainList(models []Model) []domain.User {
	result := make([]domain.User, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
