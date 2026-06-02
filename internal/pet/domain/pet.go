package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	cosmeticdomain "ditza/internal/cosmetic/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Pet struct {
	UserID               valueobject.UserID
	Name                 string
	Level                int
	XP                   int
	Mood                 Mood
	EquippedHatID        *valueobject.CosmeticID
	EquippedShirtID      *valueobject.CosmeticID
	EquippedBackgroundID *valueobject.CosmeticID
	EquippedAccessoryID  *valueobject.CosmeticID
	LastInteractionAt    time.Time
	UpdatedAt            time.Time
}

func New(userID valueobject.UserID, name string) (*Pet, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Ditza"
	}
	if utf8.RuneCountInString(name) > valueobject.MaxPetNameLength {
		return nil, domainerror.New("INVALID_PET_NAME", "el nombre de la mascota excede la longitud máxima", domainerror.ErrInvalidInput)
	}

	now := time.Now().UTC()
	return &Pet{
		UserID:            userID,
		Name:              name,
		Level:             1,
		XP:                0,
		Mood:              MoodNeutral,
		LastInteractionAt: now,
		UpdatedAt:         now,
	}, nil
}

func (p *Pet) AddXP(amount int) {
	if amount <= 0 {
		return
	}
	p.XP += amount
	for p.XP >= valueobject.PetXPPerLevel {
		p.XP -= valueobject.PetXPPerLevel
		p.Level++
	}
	p.UpdatedAt = time.Now().UTC()
}

func (p *Pet) RegisterInteraction(at time.Time) {
	p.LastInteractionAt = at
	p.UpdatedAt = time.Now().UTC()
}

func (p *Pet) UpdateMoodFromProgress(completedToday, activeHabits int) {
	if activeHabits == 0 {
		p.Mood = MoodNeutral
		p.UpdatedAt = time.Now().UTC()
		return
	}

	threshold := (activeHabits + 1) / 2
	switch {
	case completedToday >= threshold:
		p.Mood = MoodHappy
	case completedToday == 0:
		p.Mood = MoodSad
	default:
		p.Mood = MoodNeutral
	}
	p.UpdatedAt = time.Now().UTC()
}

func (p *Pet) MarkAsSleeping() {
	p.Mood = MoodSleeping
	p.UpdatedAt = time.Now().UTC()
}

func (p *Pet) EquipCosmetic(slot cosmeticdomain.Slot, cosmeticID valueobject.CosmeticID) error {
	switch slot {
	case cosmeticdomain.SlotHat:
		p.EquippedHatID = &cosmeticID
	case cosmeticdomain.SlotShirt:
		p.EquippedShirtID = &cosmeticID
	case cosmeticdomain.SlotBackground:
		p.EquippedBackgroundID = &cosmeticID
	case cosmeticdomain.SlotAccessory:
		p.EquippedAccessoryID = &cosmeticID
	default:
		return domainerror.New("INVALID_COSMETIC_SLOT", "tipo de cosmético incompatible con la ranura", domainerror.ErrInvalidCosmeticSlot)
	}
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Pet) UnequipCosmetic(slot cosmeticdomain.Slot) error {
	switch slot {
	case cosmeticdomain.SlotHat:
		p.EquippedHatID = nil
	case cosmeticdomain.SlotShirt:
		p.EquippedShirtID = nil
	case cosmeticdomain.SlotBackground:
		p.EquippedBackgroundID = nil
	case cosmeticdomain.SlotAccessory:
		p.EquippedAccessoryID = nil
	default:
		return domainerror.New("INVALID_COSMETIC_SLOT", "tipo de cosmético incompatible con la ranura", domainerror.ErrInvalidCosmeticSlot)
	}
	p.UpdatedAt = time.Now().UTC()
	return nil
}
