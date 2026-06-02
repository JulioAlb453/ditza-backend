package http

type FriendRankingEntryDTO struct {
	UserID        string `json:"user_id"`
	Alias         string `json:"alias"`
	SeasonPoints  int    `json:"season_points"`
	Rank          int    `json:"rank"`
	IsCurrentUser bool   `json:"is_current_user"`
}
