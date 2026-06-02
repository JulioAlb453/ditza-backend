package data

import (
	"ditza/internal/cosmetic/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

func ToDomain(m Model) domain.Cosmetic {
	monitoring.Mapper(logger.ModelCosmetic, "to_domain", m.ID)
	return domain.Cosmetic{
		ID:         valueobject.CosmeticID(m.ID),
		Name:       m.Name,
		Slot:       domain.Slot(m.Slot),
		PriceCoins: m.PriceCoins,
		Rarity:     domain.Rarity(m.Rarity),
		AssetKey:   m.AssetKey,
		IsActive:   m.IsActive,
		CreatedAt:  m.CreatedAt,
	}
}

func ToModel(e domain.Cosmetic) Model {
	monitoring.Mapper(logger.ModelCosmetic, "to_model", e.ID)
	return Model{
		ID:         uint64(e.ID),
		Name:       e.Name,
		Slot:       e.Slot.String(),
		PriceCoins: e.PriceCoins,
		Rarity:     e.Rarity.String(),
		AssetKey:   e.AssetKey,
		IsActive:   e.IsActive,
		CreatedAt:  e.CreatedAt,
	}
}

func ToDomainList(models []Model) []domain.Cosmetic {
	result := make([]domain.Cosmetic, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
