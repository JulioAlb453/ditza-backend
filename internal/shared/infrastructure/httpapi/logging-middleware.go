package httpapi

import (
	"net/http"
	"time"

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

		userID := r.Header.Get("X-User-ID")
		next.ServeHTTP(recorder, r)

		monitoring.HTTPRequest(
			r.Method,
			r.URL.Path,
			recorder.status,
			time.Since(startedAt),
			userID,
		)
	})
}
