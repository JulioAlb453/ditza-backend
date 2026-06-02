package http

type ActiveSeasonResponseDTO struct {
	SeasonID uint64 `json:"season_id"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	IsActive bool   `json:"is_active"`
}
