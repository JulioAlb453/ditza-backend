package domain

import (
	"time"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type PointTransaction struct {
	ID          valueobject.PointTransactionID
	UserID      valueobject.UserID
	Type        Type
	CoinsDelta  int
	SeasonDelta int
	ReferenceID *uint64
	CreatedAt   time.Time
}

func New(
	userID valueobject.UserID,
	txType Type,
	coinsDelta int,
	seasonDelta int,
	referenceID *uint64,
) (*PointTransaction, error) {
	if !txType.IsValid() {
		return nil, domainerror.New("INVALID_TRANSACTION_TYPE", "tipo de transacción de puntos inválido", domainerror.ErrInvalidInput)
	}
	if coinsDelta == 0 && seasonDelta == 0 {
		return nil, domainerror.New("INVALID_TRANSACTION_AMOUNT", "la transacción debe modificar monedas o puntos de temporada", domainerror.ErrInvalidInput)
	}

	return &PointTransaction{
		UserID:      userID,
		Type:        txType,
		CoinsDelta:  coinsDelta,
		SeasonDelta: seasonDelta,
		ReferenceID: referenceID,
		CreatedAt:   time.Now().UTC(),
	}, nil
}
