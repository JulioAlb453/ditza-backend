package domain

import "fmt"

type Rarity string

const (
	RarityCommon Rarity = "common"
	RarityRare   Rarity = "rare"
)

func ParseRarity(raw string) (Rarity, error) {
	rarity := Rarity(raw)
	switch rarity {
	case RarityCommon, RarityRare:
		return rarity, nil
	default:
		return "", fmt.Errorf("rareza de cosmético inválida: %s", raw)
	}
}

func (r Rarity) String() string {
	return string(r)
}

func (r Rarity) IsValid() bool {
	_, err := ParseRarity(string(r))
	return err == nil
}
