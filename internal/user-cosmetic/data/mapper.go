package data

import (
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/user-cosmetic/domain"
)

func ToDomain(m Model) domain.UserCosmetic {
	return domain.UserCosmetic{
		UserID:      valueobject.UserID(m.UserID),
		CosmeticID:  valueobject.CosmeticID(m.CosmeticID),
		PurchasedAt: m.PurchasedAt,
	}
}

func ToModel(e domain.UserCosmetic) Model {
	return Model{
		UserID:      uint64(e.UserID),
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
