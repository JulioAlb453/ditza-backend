package application

import (
	"context"

	leaderboarddomain "ditza/internal/leaderboard/domain"
	seasondomain "ditza/internal/season/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
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

func (s *Service) GetFriendRanking(ctx context.Context, userID valueobject.UserID) (entries []leaderboarddomain.FriendEntry, err error) {
	tracker := monitoring.StartService(logger.ModelLeaderboard, "get_friend_ranking", map[string]any{"user_id": userID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"count": len(entries)})
	}()

	activeSeason, err := s.seasonRepository.FindActive(ctx)
	if err != nil {
		return nil, err
	}
	if activeSeason == nil {
		return nil, domainerror.New("SEASON_NOT_ACTIVE", "no hay una temporada activa", domainerror.ErrSeasonNotActive)
	}
	entries, err = s.leaderboardRepository.GetFriendRanking(ctx, userID, activeSeason.ID)
	return entries, err
}
