package data

import (
	"ditza/internal/pet/domain"
	valueobject "ditza/internal/shared/domain/value-object"
)
func ToDomain(m Model) domain.Pet {
	return domain.Pet{
		UserID:               valueobject.UserID(m.UserID),
		Name:                 m.Name,
		Level:                m.Level,
		XP:                   m.XP,
		Mood:                 domain.Mood(m.Mood),
		EquippedHatID:        uint64PtrToCosmeticID(m.EquippedHatID),
		EquippedShirtID:      uint64PtrToCosmeticID(m.EquippedShirtID),
		EquippedBackgroundID: uint64PtrToCosmeticID(m.EquippedBackgroundID),
		EquippedAccessoryID:  uint64PtrToCosmeticID(m.EquippedAccessoryID),
		LastInteractionAt:    m.LastInteractionAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func ToModel(e domain.Pet) Model {
	return Model{
		UserID:               uint64(e.UserID),
		Name:                 e.Name,
		Level:                e.Level,
		XP:                   e.XP,
		Mood:                 e.Mood.String(),
		EquippedHatID:        cosmeticIDPtrToUint64(e.EquippedHatID),
		EquippedShirtID:      cosmeticIDPtrToUint64(e.EquippedShirtID),
		EquippedBackgroundID: cosmeticIDPtrToUint64(e.EquippedBackgroundID),
		EquippedAccessoryID:  cosmeticIDPtrToUint64(e.EquippedAccessoryID),
		LastInteractionAt:    e.LastInteractionAt,
		UpdatedAt:            e.UpdatedAt,
	}
}

func uint64PtrToCosmeticID(id *uint64) *valueobject.CosmeticID {
	if id == nil {
		return nil
	}
	cosmeticID := valueobject.CosmeticID(*id)
	return &cosmeticID
}

func cosmeticIDPtrToUint64(id *valueobject.CosmeticID) *uint64 {
	if id == nil {
		return nil
	}
	value := uint64(*id)
	return &value
}
