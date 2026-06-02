package http

import (
	"net/http"
	"time"

	seasonapp "ditza/internal/season/application"
	"ditza/internal/shared/infrastructure/httpapi"
)

type Controller struct {
	service *seasonapp.Service
}

func NewController(service *seasonapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetActive(w http.ResponseWriter, r *http.Request) {
	seasonEntity, err := c.service.GetActive(r.Context())
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, ActiveSeasonResponseDTO{
		SeasonID: uint64(seasonEntity.ID),
		StartsAt: seasonEntity.StartsAt.UTC().Format(time.RFC3339),
		EndsAt:   seasonEntity.EndsAt.UTC().Format(time.RFC3339),
		IsActive: seasonEntity.IsActive,
	})
}
