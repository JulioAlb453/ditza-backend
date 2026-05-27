package http

import (
	"net/http"

	leaderboardapp "ditza/internal/leaderboard/application"
	"ditza/internal/shared/infrastructure/httpapi"
)

type Controller struct {
	service *leaderboardapp.Service
}

func NewController(service *leaderboardapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetFriendRanking(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	ranking, err := c.service.GetFriendRanking(r.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	response := make([]FriendRankingEntryDTO, 0, len(ranking))
	for _, item := range ranking {
		response = append(response, FriendRankingEntryDTO{
			UserID:        uint64(item.UserID),
			Name:          item.Name,
			SeasonPoints:  item.SeasonPoints,
			Rank:          item.Rank,
			IsCurrentUser: item.IsCurrentUser,
		})
	}

	httpapi.WriteJSON(w, http.StatusOK, response)
}
