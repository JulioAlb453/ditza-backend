package domainerror

import "errors"

var (
	ErrNotFound              = errors.New("recurso no encontrado")
	ErrUnauthorized          = errors.New("no autorizado")
	ErrInvalidInput          = errors.New("datos de entrada inválidos")
	ErrHabitLimitReached     = errors.New("límite de hábitos activos alcanzado")
	ErrHabitAlreadyCompleted = errors.New("el hábito ya fue completado hoy")
	ErrHabitNotOwned         = errors.New("el hábito no pertenece al usuario")
	ErrInsufficientCoins     = errors.New("monedas insuficientes")
	ErrCosmeticAlreadyOwned  = errors.New("ya posees este cosmético")
	ErrCosmeticNotOwned      = errors.New("no posees este cosmético")
	ErrInvalidCosmeticSlot   = errors.New("tipo de cosmético incompatible con la ranura")
	ErrFriendshipExists      = errors.New("ya existe una solicitud o amistad con este usuario")
	ErrFriendshipNotPending  = errors.New("la solicitud de amistad no está pendiente")
	ErrCannotFriendSelf      = errors.New("no puedes agregarte a ti mismo como amigo")
	ErrSeasonNotActive       = errors.New("no hay una temporada activa")
	ErrInvalidTimezone       = errors.New("zona horaria inválida")
)

type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

func New(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}
