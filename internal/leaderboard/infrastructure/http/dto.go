package http

type FriendRankingEntryDTO struct {
	UserID        uint64 `json:"user_id"`
	Name          string `json:"name"`
	SeasonPoints  int    `json:"season_points"`
	Rank          int    `json:"rank"`
	IsCurrentUser bool   `json:"is_current_user"`
}
