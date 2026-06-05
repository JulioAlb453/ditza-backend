package domain

import (
	"fmt"
	"regexp"
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
	Description       string
	Category          string
	Color             string
	Frequency         string
	TargetCount       int
	TargetUnit        string
	Difficulty        string
	ReminderTime      string
	IsActive          bool
	CurrentStreak     int
	BestStreak        int
	LastCompletedDate *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type HabitConfig struct {
	Title        string
	Description  string
	Category     string
	Color        string
	Frequency    string
	TargetCount  int
	TargetUnit   string
	Difficulty   string
	ReminderTime string
}

const (
	DefaultHabitCategory    = "general"
	DefaultHabitColor       = "green"
	DefaultHabitFrequency   = "daily"
	DefaultHabitTargetCount = 1
	DefaultHabitTargetUnit  = "veces"
	DefaultHabitDifficulty  = "medium"
)

var reminderTimePattern = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

func New(userID valueobject.UserID, config HabitConfig) (*Habit, error) {
	normalized, err := normalizeConfig(config)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Habit{
		UserID:        userID,
		Title:         normalized.Title,
		Description:   normalized.Description,
		Category:      normalized.Category,
		Color:         normalized.Color,
		Frequency:     normalized.Frequency,
		TargetCount:   normalized.TargetCount,
		TargetUnit:    normalized.TargetUnit,
		Difficulty:    normalized.Difficulty,
		ReminderTime:  normalized.ReminderTime,
		IsActive:      true,
		CurrentStreak: 0,
		BestStreak:    0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func normalizeConfig(config HabitConfig) (HabitConfig, error) {
	config.Title = strings.TrimSpace(config.Title)
	if config.Title == "" {
		return HabitConfig{}, domainerror.New("INVALID_HABIT_TITLE", "el título del hábito es obligatorio", domainerror.ErrInvalidInput)
	}
	if utf8.RuneCountInString(config.Title) > valueobject.MaxHabitTitleLength {
		return HabitConfig{}, domainerror.New("INVALID_HABIT_TITLE", "el título del hábito excede la longitud máxima", domainerror.ErrInvalidInput)
	}

	config.Description = strings.TrimSpace(config.Description)
	if utf8.RuneCountInString(config.Description) > valueobject.MaxHabitDescriptionLength {
		return HabitConfig{}, domainerror.New("INVALID_HABIT_DESCRIPTION", "la descripción del hábito excede la longitud máxima", domainerror.ErrInvalidInput)
	}

	config.Category = defaultIfBlank(config.Category, DefaultHabitCategory)
	if err := validateLength("category", config.Category, valueobject.MaxHabitCategoryLength); err != nil {
		return HabitConfig{}, err
	}

	config.Color = defaultIfBlank(config.Color, DefaultHabitColor)
	if err := validateLength("color", config.Color, valueobject.MaxHabitColorLength); err != nil {
		return HabitConfig{}, err
	}

	config.Frequency = defaultIfBlank(config.Frequency, DefaultHabitFrequency)
	if !isOneOf(config.Frequency, "daily", "weekly", "specific_days") {
		return HabitConfig{}, domainerror.New("INVALID_HABIT_FREQUENCY", "la frecuencia debe ser daily, weekly o specific_days", domainerror.ErrInvalidInput)
	}

	if config.TargetCount == 0 {
		config.TargetCount = DefaultHabitTargetCount
	}
	if config.TargetCount < 0 {
		return HabitConfig{}, domainerror.New("INVALID_HABIT_TARGET", "la meta del hábito debe ser mayor a cero", domainerror.ErrInvalidInput)
	}

	config.TargetUnit = defaultIfBlank(config.TargetUnit, DefaultHabitTargetUnit)
	if err := validateLength("target_unit", config.TargetUnit, valueobject.MaxHabitTargetUnitLength); err != nil {
		return HabitConfig{}, err
	}

	config.Difficulty = defaultIfBlank(config.Difficulty, DefaultHabitDifficulty)
	if !isOneOf(config.Difficulty, "easy", "medium", "hard") {
		return HabitConfig{}, domainerror.New("INVALID_HABIT_DIFFICULTY", "la dificultad debe ser easy, medium o hard", domainerror.ErrInvalidInput)
	}

	config.ReminderTime = strings.TrimSpace(config.ReminderTime)
	if config.ReminderTime != "" && !reminderTimePattern.MatchString(config.ReminderTime) {
		return HabitConfig{}, domainerror.New("INVALID_HABIT_REMINDER_TIME", "la hora de recordatorio debe tener formato HH:MM", domainerror.ErrInvalidInput)
	}

	return config, nil
}

func defaultIfBlank(value, fallback string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return fallback
	}
	return value
}

func validateLength(field, value string, max int) error {
	if utf8.RuneCountInString(value) <= max {
		return nil
	}
	return domainerror.New(
		fmt.Sprintf("INVALID_HABIT_%s", strings.ToUpper(field)),
		fmt.Sprintf("el campo %s del hábito excede la longitud máxima", field),
		domainerror.ErrInvalidInput,
	)
}

func isOneOf(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}
	return false
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
