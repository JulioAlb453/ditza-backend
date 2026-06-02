package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Habit struct {
	ID                valueobject.HabitID
	UserID            valueobject.UserID
	Title             string
	IsActive          bool
	CurrentStreak     int
	BestStreak        int
	LastCompletedDate *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func New(userID valueobject.UserID, title string) (*Habit, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, domainerror.New("INVALID_HABIT_TITLE", "el título del hábito es obligatorio", domainerror.ErrInvalidInput)
	}
	if utf8.RuneCountInString(title) > valueobject.MaxHabitTitleLength {
		return nil, domainerror.New("INVALID_HABIT_TITLE", "el título del hábito excede la longitud máxima", domainerror.ErrInvalidInput)
	}

	now := time.Now().UTC()
	return &Habit{
		UserID:        userID,
		Title:         title,
		IsActive:      true,
		CurrentStreak: 0,
		BestStreak:    0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (h *Habit) BelongsTo(userID valueobject.UserID) bool {
	return h.UserID == userID
}

func (h *Habit) IsCompletedOn(date time.Time) bool {
	if h.LastCompletedDate == nil {
		return false
	}
	return sameCalendarDay(*h.LastCompletedDate, date)
}

func (h *Habit) CanCompleteOn(date time.Time) error {
	if !h.IsActive {
		return domainerror.New("HABIT_INACTIVE", "el hábito está inactivo", domainerror.ErrInvalidInput)
	}
	if h.IsCompletedOn(date) {
		return domainerror.New("HABIT_ALREADY_COMPLETED", "el hábito ya fue completado hoy", domainerror.ErrHabitAlreadyCompleted)
	}
	return nil
}

func (h *Habit) ApplyCompletion(completedAt time.Time) {
	h.updateStreak(completedAt)
	completedDate := truncateToDate(completedAt)
	h.LastCompletedDate = &completedDate
	h.UpdatedAt = time.Now().UTC()
}

func (h *Habit) Deactivate() {
	h.IsActive = false
	h.UpdatedAt = time.Now().UTC()
}

func (h *Habit) updateStreak(completedAt time.Time) {
	completedDate := truncateToDate(completedAt)

	if h.LastCompletedDate == nil {
		h.CurrentStreak = 1
	} else {
		yesterday := completedDate.AddDate(0, 0, -1)
		if sameCalendarDay(*h.LastCompletedDate, yesterday) {
			h.CurrentStreak++
		} else if sameCalendarDay(*h.LastCompletedDate, completedDate) {
			return
		} else {
			h.CurrentStreak = 1
		}
	}

	if h.CurrentStreak > h.BestStreak {
		h.BestStreak = h.CurrentStreak
	}
}

func sameCalendarDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func truncateToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
