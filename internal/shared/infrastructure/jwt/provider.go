package jwt

import (
	"fmt"
	"time"

	valueobject "ditza/internal/shared/domain/value-object"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwtlib.RegisteredClaims
}

type Provider struct {
	secret     []byte
	expiration time.Duration
}

func NewProvider(secret string, expiration time.Duration) (*Provider, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET debe tener al menos 32 caracteres")
	}
	if expiration <= 0 {
		expiration = 72 * time.Hour
	}
	return &Provider{
		secret:     []byte(secret),
		expiration: expiration,
	}, nil
}

func (p *Provider) Generate(userID valueobject.UserID, email string) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().UTC().Add(p.expiration)
	claims := Claims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now().UTC()),
			Subject:   userID.String(),
		},
	}

	signed, err := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString(p.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("no se pudo generar el token: %w", err)
	}
	return signed, expiresAt, nil
}

func (p *Provider) Parse(tokenString string) (valueobject.UserID, error) {
	parsed, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (any, error) {
		if token.Method != jwtlib.SigningMethodHS256 {
			return nil, fmt.Errorf("método de firma inválido")
		}
		return p.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("token inválido o expirado")
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return "", fmt.Errorf("token inválido o expirado")
	}

	userID, err := valueobject.ParseUserID(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("token inválido o expirado")
	}
	return userID, nil
}
