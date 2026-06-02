package httpapi

import (
	"fmt"
	"net/http"
	"strings"

	valueobject "ditza/internal/shared/domain/value-object"
	jwtprovider "ditza/internal/shared/infrastructure/jwt"
)

var authProvider *jwtprovider.Provider

func InitAuth(provider *jwtprovider.Provider) {
	authProvider = provider
}

func ReadUserIDFromHeader(r *http.Request) (valueobject.UserID, error) {
	if authProvider == nil {
		return "", fmt.Errorf("autenticación no configurada")
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", fmt.Errorf("falta el header Authorization con Bearer token")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("formato de Authorization inválido, use Bearer <token>")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if token == "" {
		return "", fmt.Errorf("token vacío")
	}

	return authProvider.Parse(token)
}

func TryReadUserID(r *http.Request) string {
	userID, err := ReadUserIDFromHeader(r)
	if err != nil {
		return ""
	}
	return userID.String()
}
