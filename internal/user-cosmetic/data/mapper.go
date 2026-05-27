package data

import (
	"fmt"

	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	"ditza/internal/user-cosmetic/domain"
)

func ToDomain(m Model) domain.UserCosmetic {
	monitoring.Mapper(logger.ModelUserCosmetic, "to_domain", map[string]string{
		"user_id":     m.UserID,
		"cosmetic_id": fmt.Sprintf("%d", m.CosmeticID),
	})
	return domain.UserCosmetic{
		UserID:      valueobject.UserID(m.UserID),
		CosmeticID:  valueobject.CosmeticID(m.CosmeticID),
		PurchasedAt: m.PurchasedAt,
	}
}

func ToModel(e domain.UserCosmetic) Model {
	monitoring.Mapper(logger.ModelUserCosmetic, "to_model", map[string]string{
		"user_id":     string(e.UserID),
		"cosmetic_id": fmt.Sprintf("%d", e.CosmeticID),
	})
	return Model{
		UserID:      string(e.UserID),
		CosmeticID:  uint64(e.CosmeticID),
		PurchasedAt: e.PurchasedAt,
	}
}

func ToDomainList(models []Model) []domain.UserCosmetic {
	result := make([]domain.UserCosmetic, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
