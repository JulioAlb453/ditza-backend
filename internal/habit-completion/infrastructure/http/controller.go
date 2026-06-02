package http

import (
	"net/http"
	"strconv"
	"time"

	habitcompletionapp "ditza/internal/habit-completion/application"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/httpapi"
)

type Controller struct {
	service *habitcompletionapp.Service
}

func NewController(service *habitcompletionapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Complete(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	habitIDRaw := r.PathValue("id")
	habitID, err := strconv.ParseUint(habitIDRaw, 10, 64)
	if err != nil || habitID == 0 {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_ID", Message: "id de hábito inválido"})
		return
	}

	var request CompleteHabitRequestDTO
	if err := httpapi.DecodeJSON(r, &request); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_BODY", Message: "cuerpo de petición inválido"})
		return
	}

	result, err := c.service.Complete(r.Context(), habitcompletionapp.CompleteHabitCommand{
		UserID:      userID,
		HabitID:     valueobject.HabitID(habitID),
		CompletedAt: time.Now().UTC(),
		Note:        request.Note,
		Emoji:       request.Emoji,
	})
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, CompleteHabitResponseDTO{
		CoinsEarned:         result.CoinsEarned,
		SeasonPointsEarned:  result.SeasonPointsEarned,
		CurrentStreak:       result.CurrentStreak,
		WalletCoins:         result.WalletCoins,
		CurrentSeasonPoints: result.CurrentSeasonPoints,
		PetLevel:            result.PetLevel,
		PetMood:             result.PetMood,
	})
}
