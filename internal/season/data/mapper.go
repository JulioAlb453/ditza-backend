package data

import (
	"ditza/internal/season/domain"
	valueobject "ditza/internal/shared/domain/value-object"
)
func ToDomain(m Model) domain.Season {
	return domain.Season{
		ID:        valueobject.SeasonID(m.ID),
		StartsAt:  m.StartsAt,
		EndsAt:    m.EndsAt,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
	}
}

func ToModel(e domain.Season) Model {
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
