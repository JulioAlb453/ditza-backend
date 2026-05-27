package http

type CompleteHabitRequestDTO struct {
	Note  *string `json:"note"`
	Emoji *string `json:"emoji"`
}

type CompleteHabitResponseDTO struct {
	CoinsEarned         int    `json:"coins_earned"`
	SeasonPointsEarned  int    `json:"season_points_earned"`
	CurrentStreak       int    `json:"current_streak"`
	WalletCoins         int    `json:"wallet_coins"`
	CurrentSeasonPoints int    `json:"current_season_points"`
	PetLevel            int    `json:"pet_level"`
	PetMood             string `json:"pet_mood"`
}
