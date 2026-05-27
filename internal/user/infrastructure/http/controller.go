package http

import (
	"net/http"

	"ditza/internal/shared/infrastructure/httpapi"
	userapp "ditza/internal/user/application"
)

type Controller struct {
	service *userapp.Service
}

func NewController(service *userapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequestDTO
	if err := httpapi.DecodeJSON(r, &request); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "cuerpo de petición inválido",
		})
		return
	}

	result, err := c.service.Register(r.Context(), userapp.RegisterCommand{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
		Timezone:     request.Timezone,
	})
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, RegisterResponseDTO{
		UserID:     uint64(result.UserID),
		Name:       result.Name,
		Email:      result.Email,
		Timezone:   result.Timezone,
		FriendCode: result.FriendCode,
	})
}

func (c *Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: err.Error(),
		})
		return
	}

	userEntity, err := c.service.GetByID(r.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, UserProfileResponseDTO{
		UserID:     uint64(userEntity.ID),
		Name:       userEntity.Name,
		Email:      userEntity.Email,
		Timezone:   userEntity.Timezone,
		Coins:      userEntity.Coins,
		FriendCode: userEntity.FriendCode,
	})
}
