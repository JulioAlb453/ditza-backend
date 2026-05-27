package httpapi

import (
	"fmt"
	"net/http"
	"strings"

	valueobject "ditza/internal/shared/domain/value-object"
)

const userIDHeader = "X-User-ID"

func ReadUserIDFromHeader(r *http.Request) (valueobject.UserID, error) {
	raw := strings.TrimSpace(r.Header.Get(userIDHeader))
	if raw == "" {
		return "", fmt.Errorf("falta el header %s", userIDHeader)
	}

	userID, err := valueobject.ParseUserID(raw)
	if err != nil {
		return "", fmt.Errorf("header %s inválido", userIDHeader)
	}

	return userID, nil
}
