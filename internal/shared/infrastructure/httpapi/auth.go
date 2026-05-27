package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	valueobject "ditza/internal/shared/domain/value-object"
)

const userIDHeader = "X-User-ID"

func ReadUserIDFromHeader(r *http.Request) (valueobject.UserID, error) {
	raw := strings.TrimSpace(r.Header.Get(userIDHeader))
	if raw == "" {
		return 0, fmt.Errorf("falta el header %s", userIDHeader)
	}

	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("header %s inválido", userIDHeader)
	}

	return valueobject.UserID(id), nil
}
