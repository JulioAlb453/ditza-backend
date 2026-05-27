package http

import (
	"net/http"
	"strconv"

	habitapp "ditza/internal/habit/application"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/httpapi"
)

type Controller struct {
	service *habitapp.Service
}

func NewController(service *habitapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	habits, err := c.service.ListActiveByUser(r.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	response := make([]HabitResponseDTO, 0, len(habits))
	for _, habit := range habits {
		item := HabitResponseDTO{
			HabitID:       uint64(habit.ID),
			UserID:        uint64(habit.UserID),
			Title:         habit.Title,
			IsActive:      habit.IsActive,
			CurrentStreak: habit.CurrentStreak,
			BestStreak:    habit.BestStreak,
		}
		if habit.LastCompletedDate != nil {
			item.LastCompletedDate = habit.LastCompletedDate.UTC().Format("2006-01-02")
		}
		response = append(response, item)
	}

	httpapi.WriteJSON(w, http.StatusOK, response)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	var request CreateHabitRequestDTO
	if err := httpapi.DecodeJSON(r, &request); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_BODY", Message: "cuerpo de petición inválido"})
		return
	}

	habit, err := c.service.Create(r.Context(), habitapp.CreateHabitCommand{
		UserID: userID,
		Title:  request.Title,
	})
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, HabitResponseDTO{
		HabitID:       uint64(habit.ID),
		UserID:        uint64(habit.UserID),
		Title:         habit.Title,
		IsActive:      habit.IsActive,
		CurrentStreak: habit.CurrentStreak,
		BestStreak:    habit.BestStreak,
	})
}

func (c *Controller) Deactivate(w http.ResponseWriter, r *http.Request) {
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

	if err := c.service.Deactivate(r.Context(), userID, valueobject.HabitID(habitID)); err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]string{"message": "hábito desactivado"})
}
