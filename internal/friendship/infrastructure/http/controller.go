package http

import (
	"net/http"
	"strconv"
	"time"

	friendshipapp "ditza/internal/friendship/application"
	friendshipdomain "ditza/internal/friendship/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/httpapi"
)

type Controller struct {
	service *friendshipapp.Service
}

func NewController(service *friendshipapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) SendRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	var request SendFriendRequestDTO
	if err := httpapi.DecodeJSON(r, &request); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_BODY", Message: "cuerpo de petición inválido"})
		return
	}

	addresseeID, err := valueobject.ParseUserID(request.AddresseeID)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_ID", Message: "id de usuario inválido"})
		return
	}

	friendshipEntity, err := c.service.SendRequest(r.Context(), friendshipapp.SendRequestCommand{
		RequesterID: userID,
		AddresseeID: addresseeID,
	})
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, toFriendshipDTO(friendshipEntity))
}

func (c *Controller) Accept(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	friendshipID, ok := parseFriendshipID(w, r)
	if !ok {
		return
	}

	if err := c.service.Accept(r.Context(), friendshipID, userID); err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]string{"message": "solicitud aceptada"})
}

func (c *Controller) Reject(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	friendshipID, ok := parseFriendshipID(w, r)
	if !ok {
		return
	}

	if err := c.service.Reject(r.Context(), friendshipID, userID); err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, map[string]string{"message": "solicitud rechazada"})
}

func (c *Controller) ListFriends(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	friendships, err := c.service.ListFriends(r.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	response := make([]FriendshipResponseDTO, 0, len(friendships))
	for _, item := range friendships {
		response = append(response, toFriendshipDTO(&item))
	}

	httpapi.WriteJSON(w, http.StatusOK, response)
}

func (c *Controller) ListPending(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	friendships, err := c.service.ListPendingRequests(r.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	response := make([]FriendshipResponseDTO, 0, len(friendships))
	for _, item := range friendships {
		response = append(response, toFriendshipDTO(&item))
	}

	httpapi.WriteJSON(w, http.StatusOK, response)
}

func parseFriendshipID(w http.ResponseWriter, r *http.Request) (valueobject.FriendshipID, bool) {
	raw := r.PathValue("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_ID", Message: "id de amistad inválido"})
		return 0, false
	}
	return valueobject.FriendshipID(id), true
}

func toFriendshipDTO(item *friendshipdomain.Friendship) FriendshipResponseDTO {
	response := FriendshipResponseDTO{
		FriendshipID: uint64(item.ID),
		RequesterID:  item.RequesterID.String(),
		AddresseeID:  item.AddresseeID.String(),
		Status:       item.Status.String(),
		CreatedAt:    item.CreatedAt.UTC().Format(time.RFC3339),
	}
	if item.RespondedAt != nil {
		response.RespondedAt = item.RespondedAt.UTC().Format(time.RFC3339)
	}
	return response
}
