package main

import (
	cosmeticapp "ditza/internal/cosmetic/application"
	cosmetichttp "ditza/internal/cosmetic/infrastructure/http"
	friendshipapp "ditza/internal/friendship/application"
	friendshiphttp "ditza/internal/friendship/infrastructure/http"
	habitcompletionapp "ditza/internal/habit-completion/application"
	habitcompletionhttp "ditza/internal/habit-completion/infrastructure/http"
	habitapp "ditza/internal/habit/application"
	habithttp "ditza/internal/habit/infrastructure/http"
	leaderboardapp "ditza/internal/leaderboard/application"
	leaderboardhttp "ditza/internal/leaderboard/infrastructure/http"
	seasonapp "ditza/internal/season/application"
	seasonhttp "ditza/internal/season/infrastructure/http"
	"ditza/internal/shared/infrastructure/httpserver"
	"ditza/internal/shared/infrastructure/stubrepo"
	usercosmeticapp "ditza/internal/user-cosmetic/application"
	usercosmetichttp "ditza/internal/user-cosmetic/infrastructure/http"
	userapp "ditza/internal/user/application"
	userhttp "ditza/internal/user/infrastructure/http"
)

type Container struct {
	Controllers httpserver.Controllers
}

func NewContainer() *Container {
	unitOfWork := &stubrepo.UnitOfWorkStub{}

	userRepository := &stubrepo.UserRepositoryStub{}
	habitRepository := &stubrepo.HabitRepositoryStub{}
	habitCompletionRepository := &stubrepo.HabitCompletionRepositoryStub{}
	petRepository := &stubrepo.PetRepositoryStub{}
	pointTransactionRepository := &stubrepo.PointTransactionRepositoryStub{}
	seasonRepository := &stubrepo.SeasonRepositoryStub{}
	seasonScoreRepository := &stubrepo.SeasonScoreRepositoryStub{}
	cosmeticRepository := &stubrepo.CosmeticRepositoryStub{}
	userCosmeticRepository := &stubrepo.UserCosmeticRepositoryStub{}
	friendshipRepository := &stubrepo.FriendshipRepositoryStub{}
	leaderboardRepository := &stubrepo.LeaderboardRepositoryStub{}

	userService := userapp.NewService(userRepository)
	habitService := habitapp.NewService(habitRepository)
	habitCompletionService := habitcompletionapp.NewService(
		unitOfWork,
		habitRepository,
		habitCompletionRepository,
		userRepository,
		seasonRepository,
		seasonScoreRepository,
		petRepository,
		pointTransactionRepository,
	)
	cosmeticService := cosmeticapp.NewService(cosmeticRepository)
	userCosmeticService := usercosmeticapp.NewService(
		unitOfWork,
		userRepository,
		cosmeticRepository,
		userCosmeticRepository,
		pointTransactionRepository,
	)
	friendshipService := friendshipapp.NewService(friendshipRepository)
	leaderboardService := leaderboardapp.NewService(seasonRepository, leaderboardRepository)
	seasonService := seasonapp.NewService(seasonRepository)

	return &Container{
		Controllers: httpserver.Controllers{
			User:            userhttp.NewController(userService),
			Habit:           habithttp.NewController(habitService),
			HabitCompletion: habitcompletionhttp.NewController(habitCompletionService),
			Cosmetic:        cosmetichttp.NewController(cosmeticService),
			UserCosmetic:    usercosmetichttp.NewController(userCosmeticService),
			Friendship:      friendshiphttp.NewController(friendshipService),
			Leaderboard:     leaderboardhttp.NewController(leaderboardService),
			Season:          seasonhttp.NewController(seasonService),
		},
	}
}
