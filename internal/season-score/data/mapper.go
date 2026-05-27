package data

import (
	"ditza/internal/season-score/domain"
	valueobject "ditza/internal/shared/domain/value-object"
)

func ToDomain(m Model) domain.SeasonScore {
	return domain.SeasonScore{
		UserID:    valueobject.UserID(m.UserID),
		SeasonID:  valueobject.SeasonID(m.SeasonID),
		Points:    m.Points,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToModel(e domain.SeasonScore) Model {
	return Model{
		UserID:    uint64(e.UserID),
		SeasonID:  uint64(e.SeasonID),
		Points:    e.Points,
		UpdatedAt: e.UpdatedAt,
	}
}

func ToDomainList(models []Model) []domain.SeasonScore {
	result := make([]domain.SeasonScore, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
