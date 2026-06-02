package http

import (
	"net/http"

	cosmeticapp "ditza/internal/cosmetic/application"
	"ditza/internal/shared/infrastructure/httpapi"
)

type Controller struct {
	service *cosmeticapp.Service
}

func NewController(service *cosmeticapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) ListActive(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.ListActive(r.Context())
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	response := make([]CosmeticResponseDTO, 0, len(items))
	for _, item := range items {
		response = append(response, CosmeticResponseDTO{
			CosmeticID: uint64(item.ID),
			Name:       item.Name,
			Slot:       item.Slot.String(),
			PriceCoins: item.PriceCoins,
			Rarity:     item.Rarity.String(),
			AssetKey:   item.AssetKey,
			IsActive:   item.IsActive,
		})
	}

	httpapi.WriteJSON(w, http.StatusOK, response)
}
