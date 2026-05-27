package http

type BuyCosmeticRequestDTO struct {
	CosmeticID uint64 `json:"cosmetic_id"`
}

type BuyCosmeticResponseDTO struct {
	CosmeticID  uint64 `json:"cosmetic_id"`
	WalletCoins int    `json:"wallet_coins"`
}

type InventoryItemResponseDTO struct {
	CosmeticID  uint64 `json:"cosmetic_id"`
	PurchasedAt string `json:"purchased_at"`
}
