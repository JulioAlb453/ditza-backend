package domain

import (
	petdomain "ditza/internal/pet/domain"
	seasondomain "ditza/internal/season/domain"
	seasonscoredomain "ditza/internal/season-score/domain"
	userdomain "ditza/internal/user/domain"
)

type UserProfileSummary struct {
	User         userdomain.User
	Pet          petdomain.Pet
	SeasonScore  seasonscoredomain.SeasonScore
	ActiveSeason seasondomain.Season
	HabitsToday  []HabitWithStatus
}
