package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type HabitCompletion struct {
	ID                  valueobject.HabitCompletionID
	HabitID             valueobject.HabitID
	UserID              valueobject.UserID
	CompletedAt         time.Time
	Note                *string
	Emoji               *string
	CoinsAwarded        int
	SeasonPointsAwarded int
	CreatedAt           time.Time
}

type CompletionReward struct {
	Coins        int
	SeasonPoints int
	StreakBonus  int
	NoteBonus    int
}

func CalculateCompletionReward(streakAfterCompletion int, hasNote bool) CompletionReward {
	reward := CompletionReward{
		Coins:        valueobject.BaseCompletionCoins,
		SeasonPoints: valueobject.BaseCompletionSeasonPoints,
	}

	switch {
	case streakAfterCompletion >= 7 && streakAfterCompletion%7 == 0:
		reward.StreakBonus = valueobject.StreakBonus7Days
	case streakAfterCompletion >= 3 && streakAfterCompletion%3 == 0:
		reward.StreakBonus = valueobject.StreakBonus3Days
	}

	if hasNote {
		reward.NoteBonus = valueobject.NoteBonusCoins
	}

	reward.Coins += reward.StreakBonus
	reward.SeasonPoints += reward.StreakBonus

	if hasNote {
		reward.Coins += valueobject.NoteBonusCoins
		reward.SeasonPoints += valueobject.NoteBonusSeasonPoints
	}

	return reward
}

func New(
	habitID valueobject.HabitID,
	userID valueobject.UserID,
	completedAt time.Time,
	note *string,
	emoji *string,
	reward CompletionReward,
) (*HabitCompletion, error) {
	if note != nil {
		trimmed := strings.TrimSpace(*note)
		if trimmed == "" {
			note = nil
		} else {
			if utf8.RuneCountInString(trimmed) > valueobject.MaxNoteLength {
				return nil, domainerror.New("INVALID_NOTE", "la nota excede la longitud máxima", domainerror.ErrInvalidInput)
			}
			note = &trimmed
		}
	}

	if emoji != nil {
		trimmed := strings.TrimSpace(*emoji)
		if trimmed == "" {
			emoji = nil
		} else {
			emoji = &trimmed
		}
	}

	hasNote := note != nil || emoji != nil
	if hasNote && reward.NoteBonus == 0 {
		reward.NoteBonus = valueobject.NoteBonusCoins
		reward.Coins += valueobject.NoteBonusCoins
		reward.SeasonPoints += valueobject.NoteBonusSeasonPoints
	}

	now := time.Now().UTC()
	return &HabitCompletion{
		HabitID:             habitID,
		UserID:              userID,
		CompletedAt:         completedAt,
		Note:                note,
		Emoji:               emoji,
		CoinsAwarded:        reward.Coins,
		SeasonPointsAwarded: reward.SeasonPoints,
		CreatedAt:           now,
	}, nil
}

func (c *HabitCompletion) HasJournalEntry() bool {
	return c.Note != nil || c.Emoji != nil
}
