package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	domainerror "ditza/internal/shared/domain/error"
	"ditza/internal/shared/infrastructure/logger"
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
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		logger.HTTP().Error("http_request_body_read_failed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("error", readErr.Error()),
		)
		return readErr
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))

	logRequestBody(r, body)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(payload); err != nil {
		logger.HTTP().Warn("http_request_body_decode_failed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("content_type", r.Header.Get("Content-Type")),
			slog.Int("body_size", len(body)),
			slog.String("body", sanitizeBodyForLog(r, body)),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
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
		logger.HTTP().Warn("http_domain_error",
			slog.Int("status", status),
			slog.String("code", domainErr.Code),
			slog.String("message", domainErr.Message),
			slog.String("error", domainErr.Error()),
		)
		WriteJSON(w, status, ErrorResponse{
			Code:    domainErr.Code,
			Message: domainErr.Message,
		})
		return
	}

	logger.HTTP().Error("http_internal_error", slog.String("error", err.Error()))
	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "ocurrió un error inesperado",
	})
}

const maxLoggedBodyLength = 4096

func logRequestBody(r *http.Request, body []byte) {
	level := logger.HTTP().Info
	message := "http_request_body_received"
	if len(bytes.TrimSpace(body)) == 0 {
		level = logger.HTTP().Warn
		message = "http_request_body_empty"
	}

	level(message,
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("content_type", r.Header.Get("Content-Type")),
		slog.Int("body_size", len(body)),
		slog.String("body", sanitizeBodyForLog(r, body)),
	)
}

func sanitizeBodyForLog(r *http.Request, body []byte) string {
	if isSensitivePath(r.URL.Path) {
		return "[redacted sensitive body]"
	}

	text := string(bytes.TrimSpace(body))
	if text == "" {
		return ""
	}
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	if len(text) > maxLoggedBodyLength {
		return text[:maxLoggedBodyLength] + "...[truncated]"
	}
	return text
}

func isSensitivePath(path string) bool {
	return strings.HasPrefix(path, "/auth/")
}
