package data

import (
	"fmt"

	"ditza/internal/season-score/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

func ToDomain(m Model) domain.SeasonScore {
	monitoring.Mapper(logger.ModelSeasonScore, "to_domain", map[string]string{
		"user_id":   m.UserID,
		"season_id": fmt.Sprintf("%d", m.SeasonID),
	})
	return domain.SeasonScore{
		UserID:    valueobject.UserID(m.UserID),
		SeasonID:  valueobject.SeasonID(m.SeasonID),
		Points:    m.Points,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToModel(e domain.SeasonScore) Model {
	monitoring.Mapper(logger.ModelSeasonScore, "to_model", map[string]string{
		"user_id":   string(e.UserID),
		"season_id": fmt.Sprintf("%d", e.SeasonID),
	})
	return Model{
		UserID:    string(e.UserID),
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
