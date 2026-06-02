package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	domainerror "ditza/internal/shared/domain/error"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func DecodeJSON(r *http.Request, payload any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(payload)
}

func WriteError(w http.ResponseWriter, err error) {
	var domainErr *domainerror.DomainError
	if errors.As(err, &domainErr) {
		status := http.StatusBadRequest
		switch {
		case errors.Is(domainErr, domainerror.ErrNotFound):
			status = http.StatusNotFound
		case errors.Is(domainErr, domainerror.ErrUnauthorized):
			status = http.StatusUnauthorized
		case errors.Is(domainErr, domainerror.ErrNotImplemented):
			status = http.StatusNotImplemented
		}
		WriteJSON(w, status, ErrorResponse{
			Code:    domainErr.Code,
			Message: domainErr.Message,
		})
		return
	}

	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "ocurrió un error inesperado",
	})
}
