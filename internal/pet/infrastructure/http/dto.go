package http

type PetResponseDTO struct {
	UserID                string  `json:"user_id"`
	Name                  string  `json:"name"`
	Level                 int     `json:"level"`
	XP                    int     `json:"xp"`
	Mood                  string  `json:"mood"`
	EquippedHatID         *uint64 `json:"equipped_hat_id,omitempty"`
	EquippedShirtID       *uint64 `json:"equipped_shirt_id,omitempty"`
	EquippedBackgroundID  *uint64 `json:"equipped_background_id,omitempty"`
	EquippedAccessoryID   *uint64 `json:"equipped_accessory_id,omitempty"`
	LastInteractionAt     string  `json:"last_interaction_at"`
}

type EquipCosmeticRequestDTO struct {
	CosmeticID *uint64 `json:"cosmetic_id"`
	Slot       string  `json:"slot,omitempty"`
	Unequip    bool    `json:"unequip,omitempty"`
}
