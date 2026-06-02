package monitoring

import (
	"log/slog"
	"time"

	"ditza/internal/shared/infrastructure/logger"
)

func Repository(model, operation string, attrs map[string]any) {
	fields := baseFields(model, "repository", operation, attrs)
	logger.Model(model).Info("repository_operation", fields...)
	logger.App().Info("repository_operation", fields...)
}

func Mapper(model, direction string, recordID any) {
	fields := []any{
		slog.String("layer", "data"),
		slog.String("model", model),
		slog.String("event", "mapper"),
		slog.String("direction", direction),
		slog.Any("record_id", recordID),
	}
	logger.Model(model).Info("mapper_operation", fields...)
}

func Service(model, operation, status string, attrs map[string]any, err error, duration time.Duration) {
	fields := baseFields(model, "application", operation, attrs)
	fields = append(fields,
		slog.String("status", status),
		slog.Int64("duration_ms", duration.Milliseconds()),
	)
	if err != nil {
		fields = append(fields, slog.String("error", err.Error()))
	}
	logger.Model(model).Info("service_operation", fields...)
	logger.App().Info("service_operation", fields...)
}

func HTTPRequest(method, path string, status int, duration time.Duration, userID string) {
	logger.HTTP().Info("http_request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", status),
		slog.Int64("duration_ms", duration.Milliseconds()),
		slog.String("user_id", userID),
	)
	logger.App().Info("http_request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", status),
	)
}

func baseFields(model, layer, operation string, attrs map[string]any) []any {
	fields := []any{
		slog.String("layer", layer),
		slog.String("model", model),
		slog.String("operation", operation),
	}
	for key, value := range attrs {
		fields = append(fields, slog.Any(key, value))
	}
	return fields
}
