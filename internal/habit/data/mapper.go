package data

import (
	"ditza/internal/habit/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

func ToDomain(m Model) domain.Habit {
	monitoring.Mapper(logger.ModelHabit, "to_domain", m.ID)
	return domain.Habit{
		ID:                valueobject.HabitID(m.ID),
		UserID:            valueobject.UserID(m.UserID),
		Title:             m.Title,
		IsActive:          m.IsActive,
		CurrentStreak:     m.CurrentStreak,
		BestStreak:        m.BestStreak,
		LastCompletedDate: m.LastCompletedDate,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func ToModel(e domain.Habit) Model {
	monitoring.Mapper(logger.ModelHabit, "to_model", e.ID)
	return Model{
		ID:                uint64(e.ID),
		UserID:            string(e.UserID),
		Title:             e.Title,
		IsActive:          e.IsActive,
		CurrentStreak:     e.CurrentStreak,
		BestStreak:        e.BestStreak,
		LastCompletedDate: e.LastCompletedDate,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

func ToDomainList(models []Model) []domain.Habit {
	result := make([]domain.Habit, len(models))
	for i, m := range models {
		result[i] = ToDomain(m)
	}
	return result
}
