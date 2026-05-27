package application

import (
	"context"

	seasondomain "ditza/internal/season/domain"
	domainerror "ditza/internal/shared/domain/error"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

type Service struct {
	seasonRepository seasondomain.Repository
}

func NewService(seasonRepository seasondomain.Repository) *Service {
	return &Service{seasonRepository: seasonRepository}
}

func (s *Service) GetActive(ctx context.Context) (activeSeason *seasondomain.Season, err error) {
	tracker := monitoring.StartService(logger.ModelSeason, "get_active", nil)
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		if activeSeason != nil {
			tracker.Success(map[string]any{"season_id": activeSeason.ID})
			return
		}
		tracker.Success(nil)
	}()

	activeSeason, err = s.seasonRepository.FindActive(ctx)
	if err != nil {
		return nil, err
	}
	if activeSeason == nil {
		return nil, domainerror.New("SEASON_NOT_ACTIVE", "no hay una temporada activa", domainerror.ErrSeasonNotActive)
	}
	return activeSeason, nil
}
