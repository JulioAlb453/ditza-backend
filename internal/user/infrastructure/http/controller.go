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
		Alias:    request.Alias,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, RegisterResponseDTO{
		UserID: result.UserID.String(),
		Alias:  result.Alias,
		Email:  result.Email,
	})
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequestDTO
	if err := httpapi.DecodeJSON(r, &request); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{
			Code:    "INVALID_BODY",
			Message: "cuerpo de petición inválido",
		})
		return
	}

	result, err := c.service.Login(r.Context(), userapp.LoginCommand{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, LoginResponseDTO{
		UserID: result.UserID.String(),
		Alias:  result.Alias,
		Email:  result.Email,
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
		UserID: userEntity.ID.String(),
		Alias:  userEntity.Alias,
		Email:  userEntity.Email,
		Coins:  userEntity.Coins,
	})
}
