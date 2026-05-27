package application

import (
	"context"

	leaderboarddomain "ditza/internal/leaderboard/domain"
	seasondomain "ditza/internal/season/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Service struct {
	seasonRepository      seasondomain.Repository
	leaderboardRepository leaderboarddomain.Repository
}

func NewService(
	seasonRepository seasondomain.Repository,
	leaderboardRepository leaderboarddomain.Repository,
) *Service {
	return &Service{
		seasonRepository:      seasonRepository,
		leaderboardRepository: leaderboardRepository,
	}
}

func (s *Service) GetFriendRanking(ctx context.Context, userID valueobject.UserID) ([]leaderboarddomain.FriendEntry, error) {
	activeSeason, err := s.seasonRepository.FindActive(ctx)
	if err != nil {
		return nil, err
	}
	if activeSeason == nil {
		return nil, domainerror.New("SEASON_NOT_ACTIVE", "no hay una temporada activa", domainerror.ErrSeasonNotActive)
	}
	return s.leaderboardRepository.GetFriendRanking(ctx, userID, activeSeason.ID)
}
