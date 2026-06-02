package domain

import (
	petdomain "ditza/internal/pet/domain"
	seasonscoredomain "ditza/internal/season-score/domain"
	seasondomain "ditza/internal/season/domain"
	userdomain "ditza/internal/user/domain"
)

type UserProfileSummary struct {
	User         userdomain.User
	Pet          petdomain.Pet
	SeasonScore  seasonscoredomain.SeasonScore
	ActiveSeason seasondomain.Season
	HabitsToday  []HabitWithStatus
}
