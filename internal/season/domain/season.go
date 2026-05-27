package domain

import (
	"time"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Season struct {
	ID        valueobject.SeasonID
	StartsAt  time.Time
	EndsAt    time.Time
	IsActive  bool
	CreatedAt time.Time
}

func New(startsAt, endsAt time.Time) (*Season, error) {
	if !endsAt.After(startsAt) {
		return nil, domainerror.New("INVALID_SEASON_RANGE", "la fecha de fin debe ser posterior al inicio", domainerror.ErrInvalidInput)
	}

	return &Season{
		StartsAt:  startsAt.UTC(),
		EndsAt:    endsAt.UTC(),
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *Season) Contains(moment time.Time) bool {
	moment = moment.UTC()
	return !moment.Before(s.StartsAt) && moment.Before(s.EndsAt)
}

func (s *Season) IsExpired(at time.Time) bool {
	return !at.UTC().Before(s.EndsAt)
}

func (s *Season) DaysRemaining(at time.Time) int {
	at = at.UTC()
	if at.After(s.EndsAt) || at.Equal(s.EndsAt) {
		return 0
	}
	remaining := s.EndsAt.Sub(at)
	days := int(remaining.Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

func (s *Season) Deactivate() {
	s.IsActive = false
}
