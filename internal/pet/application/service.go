package application

import (
	"context"

	cosmeticdomain "ditza/internal/cosmetic/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	petdomain "ditza/internal/pet/domain"
	usercosmeticdomain "ditza/internal/user-cosmetic/domain"
)

type Service struct {
	petRepository          petdomain.Repository
	cosmeticRepository     cosmeticdomain.Repository
	userCosmeticRepository usercosmeticdomain.Repository
}

type EquipCosmeticCommand struct {
	UserID     valueobject.UserID
	CosmeticID *valueobject.CosmeticID
	Slot       cosmeticdomain.Slot
	Unequip    bool
}

func NewService(
	petRepository petdomain.Repository,
	cosmeticRepository cosmeticdomain.Repository,
	userCosmeticRepository usercosmeticdomain.Repository,
) *Service {
	return &Service{
		petRepository:          petRepository,
		cosmeticRepository:     cosmeticRepository,
		userCosmeticRepository: userCosmeticRepository,
	}
}

func (s *Service) GetByUserID(ctx context.Context, userID valueobject.UserID) (petEntity *petdomain.Pet, err error) {
	tracker := monitoring.StartService(logger.ModelPet, "get_by_user_id", map[string]any{"user_id": userID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"level": petEntity.Level})
	}()

	petEntity, err = s.petRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if petEntity == nil {
		return nil, domainerror.New("PET_NOT_FOUND", "aún no tienes mascota; completa un hábito para obtenerla", domainerror.ErrNotFound)
	}
	return petEntity, nil
}

func (s *Service) EquipCosmetic(ctx context.Context, command EquipCosmeticCommand) (petEntity *petdomain.Pet, err error) {
	tracker := monitoring.StartService(logger.ModelPet, "equip_cosmetic", map[string]any{
		"user_id": command.UserID,
		"unequip": command.Unequip,
	})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(nil)
	}()

	petEntity, err = s.petRepository.FindByUserID(ctx, command.UserID)
	if err != nil {
		return nil, err
	}
	if petEntity == nil {
		return nil, domainerror.New("PET_NOT_FOUND", "aún no tienes mascota; completa un hábito para obtenerla", domainerror.ErrNotFound)
	}

	if command.Unequip {
		if !command.Slot.IsValid() {
			return nil, domainerror.New("INVALID_COSMETIC_SLOT", "tipo de cosmético incompatible con la ranura", domainerror.ErrInvalidCosmeticSlot)
		}
		if err := petEntity.UnequipCosmetic(command.Slot); err != nil {
			return nil, err
		}
	} else {
		if command.CosmeticID == nil {
			return nil, domainerror.New("INVALID_INPUT", "cosmetic_id es obligatorio para equipar", domainerror.ErrInvalidInput)
		}

		cosmeticEntity, err := s.cosmeticRepository.FindByID(ctx, *command.CosmeticID)
		if err != nil {
			return nil, err
		}
		if cosmeticEntity == nil {
			return nil, domainerror.New("COSMETIC_NOT_FOUND", "cosmético no encontrado", domainerror.ErrNotFound)
		}

		owned, err := s.userCosmeticRepository.Exists(ctx, command.UserID, *command.CosmeticID)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, domainerror.New("COSMETIC_NOT_OWNED", "no posees este cosmético", domainerror.ErrCosmeticNotOwned)
		}

		if err := petEntity.EquipCosmetic(cosmeticEntity.Slot, *command.CosmeticID); err != nil {
			return nil, err
		}
	}

	if err := s.petRepository.Update(ctx, petEntity); err != nil {
		return nil, err
	}
	return petEntity, nil
}
