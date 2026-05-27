package main

import (
	"database/sql"

	cosmeticapp "ditza/internal/cosmetic/application"
	cosmetichttp "ditza/internal/cosmetic/infrastructure/http"
	cosmeticpostgres "ditza/internal/cosmetic/infrastructure/postgres"
	friendshipapp "ditza/internal/friendship/application"
	friendshiphttp "ditza/internal/friendship/infrastructure/http"
	friendshippostgres "ditza/internal/friendship/infrastructure/postgres"
	habitcompletionapp "ditza/internal/habit-completion/application"
	habitcompletionhttp "ditza/internal/habit-completion/infrastructure/http"
	habitcompletionpostgres "ditza/internal/habit-completion/infrastructure/postgres"
	habitapp "ditza/internal/habit/application"
	habithttp "ditza/internal/habit/infrastructure/http"
	habitpostgres "ditza/internal/habit/infrastructure/postgres"
	leaderboardapp "ditza/internal/leaderboard/application"
	leaderboardhttp "ditza/internal/leaderboard/infrastructure/http"
	leaderboardpostgres "ditza/internal/leaderboard/infrastructure/postgres"
	petpostgres "ditza/internal/pet/infrastructure/postgres"
	pointtransactionpostgres "ditza/internal/point-transaction/infrastructure/postgres"
	seasonapp "ditza/internal/season/application"
	seasonhttp "ditza/internal/season/infrastructure/http"
	seasonpostgres "ditza/internal/season/infrastructure/postgres"
	seasonscorepostgres "ditza/internal/season-score/infrastructure/postgres"
	"ditza/internal/shared/infrastructure/httpserver"
	jwtprovider "ditza/internal/shared/infrastructure/jwt"
	sharedpostgres "ditza/internal/shared/infrastructure/postgres"
	usercosmeticapp "ditza/internal/user-cosmetic/application"
	usercosmetichttp "ditza/internal/user-cosmetic/infrastructure/http"
	usercosmeticpostgres "ditza/internal/user-cosmetic/infrastructure/postgres"
	userapp "ditza/internal/user/application"
	userhttp "ditza/internal/user/infrastructure/http"
	userpostgres "ditza/internal/user/infrastructure/postgres"
)

type Container struct {
	DB          *sql.DB
	Controllers httpserver.Controllers
}

func NewContainer(db *sql.DB, tokenProvider *jwtprovider.Provider) *Container {
	unitOfWork := sharedpostgres.NewUnitOfWork(db)

	userRepository := userpostgres.New(db)
	habitRepository := habitpostgres.New(db)
	habitCompletionRepository := habitcompletionpostgres.New(db)
	petRepository := petpostgres.New(db)
	pointTransactionRepository := pointtransactionpostgres.New(db)
	seasonRepository := seasonpostgres.New(db)
	seasonScoreRepository := seasonscorepostgres.New(db)
	cosmeticRepository := cosmeticpostgres.New(db)
	userCosmeticRepository := usercosmeticpostgres.New(db)
	friendshipRepository := friendshippostgres.New(db)
	leaderboardRepository := leaderboardpostgres.New(db)

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
		DB: db,
		Controllers: httpserver.Controllers{
			User:            userhttp.NewController(userService, tokenProvider),
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
