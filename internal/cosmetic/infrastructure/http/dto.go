package http

type CosmeticResponseDTO struct {
	CosmeticID uint64 `json:"cosmetic_id"`
	Name       string `json:"name"`
	Slot       string `json:"slot"`
	PriceCoins int    `json:"price_coins"`
	Rarity     string `json:"rarity"`
	AssetKey   string `json:"asset_key"`
	IsActive   bool   `json:"is_active"`
}
