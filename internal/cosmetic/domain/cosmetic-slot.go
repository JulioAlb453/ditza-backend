package domain

import "fmt"

type Slot string

const (
	SlotHat        Slot = "hat"
	SlotShirt      Slot = "shirt"
	SlotBackground Slot = "background"
	SlotAccessory  Slot = "accessory"
)

func ParseSlot(raw string) (Slot, error) {
	slot := Slot(raw)
	switch slot {
	case SlotHat, SlotShirt, SlotBackground, SlotAccessory:
		return slot, nil
	default:
		return "", fmt.Errorf("tipo de cosmético inválido: %s", raw)
	}
}

func (s Slot) String() string {
	return string(s)
}

func (s Slot) IsValid() bool {
	_, err := ParseSlot(string(s))
	return err == nil
}
