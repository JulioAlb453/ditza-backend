package http

type CreateHabitRequestDTO struct {
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	Category     string `json:"category,omitempty"`
	Color        string `json:"color,omitempty"`
	Frequency    string `json:"frequency,omitempty"`
	TargetCount  int    `json:"target_count,omitempty"`
	TargetUnit   string `json:"target_unit,omitempty"`
	Difficulty   string `json:"difficulty,omitempty"`
	ReminderTime string `json:"reminder_time,omitempty"`
}

type HabitResponseDTO struct {
	HabitID           uint64 `json:"habit_id"`
	UserID            string `json:"user_id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Category          string `json:"category"`
	Color             string `json:"color"`
	Frequency         string `json:"frequency"`
	TargetCount       int    `json:"target_count"`
	TargetUnit        string `json:"target_unit"`
	Difficulty        string `json:"difficulty"`
	ReminderTime      string `json:"reminder_time,omitempty"`
	IsActive          bool   `json:"is_active"`
	CurrentStreak     int    `json:"current_streak"`
	BestStreak        int    `json:"best_streak"`
	LastCompletedDate string `json:"last_completed_date,omitempty"`
}
