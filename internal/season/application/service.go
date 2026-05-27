package application

import (
	"context"

	seasondomain "ditza/internal/season/domain"
	domainerror "ditza/internal/shared/domain/error"
)

type Service struct {
	seasonRepository seasondomain.Repository
}

func NewService(seasonRepository seasondomain.Repository) *Service {
	return &Service{seasonRepository: seasonRepository}
}

func (s *Service) GetActive(ctx context.Context) (*seasondomain.Season, error) {
	activeSeason, err := s.seasonRepository.FindActive(ctx)
	if err != nil {
		return nil, err
	}
	if activeSeason == nil {
		return nil, domainerror.New("SEASON_NOT_ACTIVE", "no hay una temporada activa", domainerror.ErrSeasonNotActive)
	}
	return activeSeason, nil
}
