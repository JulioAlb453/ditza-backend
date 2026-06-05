package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		userID := TryReadUserID(r)
		next.ServeHTTP(recorder, r)

		duration := time.Since(startedAt)
		logHTTPStatus(r, recorder.status, duration, userID)
		monitoring.HTTPRequest(
			r.Method,
			r.URL.Path,
			recorder.status,
			duration,
			userID,
		)
	})
}

func logHTTPStatus(r *http.Request, status int, duration time.Duration, userID string) {
	fields := []any{
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", status),
		slog.Int64("duration_ms", duration.Milliseconds()),
		slog.String("user_id", userID),
	}

	switch {
	case status >= http.StatusInternalServerError:
		logger.HTTP().Error("http_response_error", fields...)
	case status >= http.StatusBadRequest:
		logger.HTTP().Warn("http_response_warning", fields...)
	}
}
