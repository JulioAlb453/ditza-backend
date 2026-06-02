package http

type CreateHabitRequestDTO struct {
	Title string `json:"title"`
}

type HabitResponseDTO struct {
	HabitID           uint64 `json:"habit_id"`
	UserID            string `json:"user_id"`
	Title             string `json:"title"`
	IsActive          bool   `json:"is_active"`
	CurrentStreak     int    `json:"current_streak"`
	BestStreak        int    `json:"best_streak"`
	LastCompletedDate string `json:"last_completed_date,omitempty"`
}
