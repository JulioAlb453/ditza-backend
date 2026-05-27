package data

import (
	"ditza/internal/habit/domain"
	valueobject "ditza/internal/shared/domain/value-object"
)
func ToDomain(m Model) domain.Habit {
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
	return Model{
		ID:                uint64(e.ID),
		UserID:            uint64(e.UserID),
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
