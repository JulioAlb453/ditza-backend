package password

import (
	"fmt"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func Hash(plain string) (string, error) {
	if len(plain) < valueobject.MinPasswordLength {
		return "", domainerror.New("INVALID_PASSWORD", "la contraseña debe tener al menos 8 caracteres", domainerror.ErrInvalidInput)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("no se pudo hashear la contraseña: %w", err)
	}
	return string(hashed), nil
}

func Verify(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
