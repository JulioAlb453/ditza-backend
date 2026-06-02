package domain

import habitdomain "ditza/internal/habit/domain"

type HabitWithStatus struct {
	Habit          habitdomain.Habit
	CompletedToday bool
}
