package domain

import (
	"strings"
	"time"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Cosmetic struct {
	ID         valueobject.CosmeticID
	Name       string
	Slot       Slot
	PriceCoins int
	Rarity     Rarity
	AssetKey   string
	IsActive   bool
	CreatedAt  time.Time
}

func New(
	name string,
	slot Slot,
	priceCoins int,
	rarity Rarity,
	assetKey string,
) (*Cosmetic, error) {
	name = strings.TrimSpace(name)
	assetKey = strings.TrimSpace(assetKey)

	if name == "" {
		return nil, domainerror.New("INVALID_COSMETIC_NAME", "el nombre del cosmético es obligatorio", domainerror.ErrInvalidInput)
	}
	if !slot.IsValid() {
		return nil, domainerror.New("INVALID_COSMETIC_SLOT", "tipo de cosmético incompatible con la ranura", domainerror.ErrInvalidCosmeticSlot)
	}
	if priceCoins <= 0 {
		return nil, domainerror.New("INVALID_COSMETIC_PRICE", "el precio del cosmético debe ser mayor a cero", domainerror.ErrInvalidInput)
	}
	if !rarity.IsValid() {
		return nil, domainerror.New("INVALID_COSMETIC_RARITY", "rareza de cosmético inválida", domainerror.ErrInvalidInput)
	}
	if assetKey == "" {
		return nil, domainerror.New("INVALID_COSMETIC_ASSET", "el identificador del asset es obligatorio", domainerror.ErrInvalidInput)
	}

	return &Cosmetic{
		Name:       name,
		Slot:       slot,
		PriceCoins: priceCoins,
		Rarity:     rarity,
		AssetKey:   assetKey,
		IsActive:   true,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

func (c *Cosmetic) Deactivate() {
	c.IsActive = false
}

func (c *Cosmetic) CanBePurchased() bool {
	return c.IsActive
}
