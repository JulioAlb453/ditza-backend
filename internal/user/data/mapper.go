package data

import (
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	"ditza/internal/user/domain"
)

func ToDomain(m Model) domain.User {
	monitoring.Mapper(logger.ModelUser, "to_domain", m.ID)
	return domain.User{
		ID:        valueobject.UserID(m.ID),
		Alias:     m.Alias,
		Email:     m.Email,
		Password:  m.Password,
		Coins:     m.Coins,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToModel(e domain.User) Model {
	monitoring.Mapper(logger.ModelUser, "to_model", e.ID)
	return Model{
		ID:        e.ID.String(),
		Alias:     e.Alias,
		Email:     e.Email,
		Password:  e.Password,
		Coins:     e.Coins,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func ToDomainList(models []Model) []domain.User {
	result := make([]domain.User, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
