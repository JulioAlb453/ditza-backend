package http

import (
	"net/http"
	"time"

	cosmeticdomain "ditza/internal/cosmetic/domain"
	petapp "ditza/internal/pet/application"
	petdomain "ditza/internal/pet/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/httpapi"
)

type Controller struct {
	service *petapp.Service
}

func NewController(service *petapp.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	petEntity, err := c.service.GetByUserID(r.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, toPetResponse(petEntity))
}

func (c *Controller) Equip(w http.ResponseWriter, r *http.Request) {
	userID, err := httpapi.ReadUserIDFromHeader(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusUnauthorized, httpapi.ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()})
		return
	}

	var request EquipCosmeticRequestDTO
	if err := httpapi.DecodeJSON(r, &request); err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_BODY", Message: "cuerpo de petición inválido"})
		return
	}

	command := petapp.EquipCosmeticCommand{
		UserID:  userID,
		Unequip: request.Unequip,
	}
	if request.CosmeticID != nil && *request.CosmeticID > 0 {
		cosmeticID := valueobject.CosmeticID(*request.CosmeticID)
		command.CosmeticID = &cosmeticID
	}
	if request.Slot != "" {
		slot, err := cosmeticdomain.ParseSlot(request.Slot)
		if err != nil {
			httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorResponse{Code: "INVALID_SLOT", Message: "tipo de cosmético inválido"})
			return
		}
		command.Slot = slot
	}

	petEntity, err := c.service.EquipCosmetic(r.Context(), command)
	if err != nil {
		httpapi.WriteError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, toPetResponse(petEntity))
}

func toPetResponse(petEntity *petdomain.Pet) PetResponseDTO {
	return PetResponseDTO{
		UserID:               petEntity.UserID.String(),
		Name:                 petEntity.Name,
		Level:                petEntity.Level,
		XP:                   petEntity.XP,
		Mood:                 petEntity.Mood.String(),
		EquippedHatID:        cosmeticIDToUint64Ptr(petEntity.EquippedHatID),
		EquippedShirtID:      cosmeticIDToUint64Ptr(petEntity.EquippedShirtID),
		EquippedBackgroundID: cosmeticIDToUint64Ptr(petEntity.EquippedBackgroundID),
		EquippedAccessoryID:  cosmeticIDToUint64Ptr(petEntity.EquippedAccessoryID),
		LastInteractionAt:    petEntity.LastInteractionAt.UTC().Format(time.RFC3339),
	}
}

func cosmeticIDToUint64Ptr(id *valueobject.CosmeticID) *uint64 {
	if id == nil {
		return nil
	}
	value := uint64(*id)
	return &value
}
