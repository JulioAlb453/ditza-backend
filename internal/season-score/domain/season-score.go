package domain

import (
	"time"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type SeasonScore struct {
	UserID    valueobject.UserID
	SeasonID  valueobject.SeasonID
	Points    int
	UpdatedAt time.Time
}

func New(userID valueobject.UserID, seasonID valueobject.SeasonID) *SeasonScore {
	return &SeasonScore{
		UserID:    userID,
		SeasonID:  seasonID,
		Points:    0,
		UpdatedAt: time.Now().UTC(),
	}
}

func (s *SeasonScore) AddPoints(amount int) error {
	if amount < 0 {
		return domainerror.New("INVALID_SEASON_POINTS", "no se pueden agregar puntos de temporada negativos", domainerror.ErrInvalidInput)
	}
	s.Points += amount
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *SeasonScore) Reset() {
	s.Points = 0
	s.UpdatedAt = time.Now().UTC()
}
