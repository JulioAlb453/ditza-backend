package data

import (
	"ditza/internal/point-transaction/domain"
	valueobject "ditza/internal/shared/domain/value-object"
)

func ToDomain(m Model) domain.PointTransaction {
	return domain.PointTransaction{
		ID:          valueobject.PointTransactionID(m.ID),
		UserID:      valueobject.UserID(m.UserID),
		Type:        domain.Type(m.Type),
		CoinsDelta:  m.CoinsDelta,
		SeasonDelta: m.SeasonDelta,
		ReferenceID: m.ReferenceID,
		CreatedAt:   m.CreatedAt,
	}
}

func ToModel(e domain.PointTransaction) Model {
	return Model{
		ID:          uint64(e.ID),
		UserID:      uint64(e.UserID),
		Type:        e.Type.String(),
		CoinsDelta:  e.CoinsDelta,
		SeasonDelta: e.SeasonDelta,
		ReferenceID: e.ReferenceID,
		CreatedAt:   e.CreatedAt,
	}
}

func ToDomainList(models []Model) []domain.PointTransaction {
	result := make([]domain.PointTransaction, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
