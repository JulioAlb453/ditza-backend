package data

import (
	"ditza/internal/friendship/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

func ToDomain(m Model) domain.Friendship {
	monitoring.Mapper(logger.ModelFriendship, "to_domain", m.ID)
	return domain.Friendship{
		ID:          valueobject.FriendshipID(m.ID),
		RequesterID: valueobject.UserID(m.RequesterID),
		AddresseeID: valueobject.UserID(m.AddresseeID),
		Status:      domain.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		RespondedAt: m.RespondedAt,
	}
}

func ToModel(e domain.Friendship) Model {
	monitoring.Mapper(logger.ModelFriendship, "to_model", e.ID)
	return Model{
		ID:          uint64(e.ID),
		RequesterID: string(e.RequesterID),
		AddresseeID: string(e.AddresseeID),
		Status:      e.Status.String(),
		CreatedAt:   e.CreatedAt,
		RespondedAt: e.RespondedAt,
	}
}

func ToDomainList(models []Model) []domain.Friendship {
	result := make([]domain.Friendship, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
