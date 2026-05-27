package data

import (
	"ditza/internal/season/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

func ToDomain(m Model) domain.Season {
	monitoring.Mapper(logger.ModelSeason, "to_domain", m.ID)
	return domain.Season{
		ID:        valueobject.SeasonID(m.ID),
		StartsAt:  m.StartsAt,
		EndsAt:    m.EndsAt,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
	}
}

func ToModel(e domain.Season) Model {
	monitoring.Mapper(logger.ModelSeason, "to_model", e.ID)
	return Model{
		ID:        uint64(e.ID),
		StartsAt:  e.StartsAt,
		EndsAt:    e.EndsAt,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
	}
}

func ToDomainList(models []Model) []domain.Season {
	result := make([]domain.Season, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
