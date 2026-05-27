package data

import (
	"ditza/internal/point-transaction/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

func ToDomain(m Model) domain.PointTransaction {
	monitoring.Mapper(logger.ModelPointTransaction, "to_domain", m.ID)
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
	monitoring.Mapper(logger.ModelPointTransaction, "to_model", e.ID)
	return Model{
		ID:          uint64(e.ID),
		UserID:      string(e.UserID),
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
