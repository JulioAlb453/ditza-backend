package data

import (
	"ditza/internal/habit-completion/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

func ToDomain(m Model) domain.HabitCompletion {
	monitoring.Mapper(logger.ModelHabitCompletion, "to_domain", m.ID)
	return domain.HabitCompletion{
		ID:                  valueobject.HabitCompletionID(m.ID),
		HabitID:             valueobject.HabitID(m.HabitID),
		UserID:              valueobject.UserID(m.UserID),
		CompletedAt:         m.CompletedAt,
		Note:                m.Note,
		Emoji:               m.Emoji,
		CoinsAwarded:        m.CoinsAwarded,
		SeasonPointsAwarded: m.SeasonPointsAwarded,
		CreatedAt:           m.CreatedAt,
	}
}

func ToModel(e domain.HabitCompletion) Model {
	monitoring.Mapper(logger.ModelHabitCompletion, "to_model", e.ID)
	return Model{
		ID:                  uint64(e.ID),
		HabitID:             uint64(e.HabitID),
		UserID:              string(e.UserID),
		CompletedAt:         e.CompletedAt,
		Note:                e.Note,
		Emoji:               e.Emoji,
		CoinsAwarded:        e.CoinsAwarded,
		SeasonPointsAwarded: e.SeasonPointsAwarded,
		CreatedAt:           e.CreatedAt,
	}
}

func ToDomainList(models []Model) []domain.HabitCompletion {
	result := make([]domain.HabitCompletion, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
