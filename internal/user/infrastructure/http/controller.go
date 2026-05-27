package http

import (
	"net/http"
	"time"

	"ditza/internal/shared/infrastructure/httpapi"
	jwtprovider "ditza/internal/shared/infrastructure/jwt"
	valueobject "ditza/internal/shared/domain/value-object"
	userapp "ditza/internal/user/application"
)

type Controller struct {
	service      *userapp.Service
	tokenProvider *jwtprovider.Provider
}

func NewController(service *userapp.Service, tokenProvider *jwtprovider.Provider) *Controller {
	return &Controller{
		service:      service,
		tokenProvider: tokenProvider,
	}
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

	response, err := c.buildAuthResponse(result.UserID, result.Alias, result.Email)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusInternalServerError, httpapi.ErrorResponse{
			Code:    "TOKEN_GENERATION_FAILED",
			Message: "no se pudo generar el token de acceso",
		})
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, response)
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

	response, err := c.buildAuthResponse(result.UserID, result.Alias, result.Email)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusInternalServerError, httpapi.ErrorResponse{
			Code:    "TOKEN_GENERATION_FAILED",
			Message: "no se pudo generar el token de acceso",
		})
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, response)
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

func (c *Controller) buildAuthResponse(userID valueobject.UserID, alias, email string) (AuthResponseDTO, error) {
	token, expiresAt, err := c.tokenProvider.Generate(userID, email)
	if err != nil {
		return AuthResponseDTO{}, err
	}

	return AuthResponseDTO{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
		UserID:      userID.String(),
		Alias:       alias,
		Email:       email,
	}, nil
}
