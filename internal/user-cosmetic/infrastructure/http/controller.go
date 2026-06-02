package http

import (
	"net/http"
	"time"

	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/httpapi"
	usercosmeticapp "ditza/internal/user-cosmetic/application"
)

type Controller struct {
	service *usercosmeticapp.Service
}

func NewController(service *usercosmeticapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Buy(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	var request BuyCosmeticRequestDTO
	if err := httpapi.DecodeJSON(r, &request); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_BODY", Message: "cuerpo de petición inválido"})
		return
	}

	result, err := c.service.Buy(r.Context(), usercosmeticapp.BuyCosmeticCommand{
		UserID:     userID,
		CosmeticID: valueobject.CosmeticID(request.CosmeticID),
	})
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, BuyCosmeticResponseDTO{
		CosmeticID:  uint64(result.CosmeticID),
		WalletCoins: result.WalletCoins,
	})
}

func (c *Controller) ListInventory(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	items, err := c.service.ListInventory(r.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	response := make([]InventoryItemResponseDTO, 0, len(items))
	for _, item := range items {
		response = append(response, InventoryItemResponseDTO{
			CosmeticID:  uint64(item.CosmeticID),
			PurchasedAt: item.PurchasedAt.UTC().Format(time.RFC3339),
		})
	}

	httpapi.WriteJSON(w, http.StatusOK, response)
}
