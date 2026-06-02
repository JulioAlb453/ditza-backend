package application

import (
	"context"

	cosmeticdomain "ditza/internal/cosmetic/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

type Service struct {
	cosmeticRepository cosmeticdomain.Repository
}

func NewService(cosmeticRepository cosmeticdomain.Repository) *Service {
	return &Service{cosmeticRepository: cosmeticRepository}
}

func (s *Service) ListActive(ctx context.Context) (items []cosmeticdomain.Cosmetic, err error) {
	tracker := monitoring.StartService(logger.ModelCosmetic, "list_active", nil)
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"count": len(items)})
	}()

	items, err = s.cosmeticRepository.FindAllActive(ctx)
	return items, err
}

func (s *Service) GetByID(ctx context.Context, cosmeticID valueobject.CosmeticID) (cosmeticEntity *cosmeticdomain.Cosmetic, err error) {
	tracker := monitoring.StartService(logger.ModelCosmetic, "get_by_id", map[string]any{"cosmetic_id": cosmeticID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(nil)
	}()

	cosmeticEntity, err = s.cosmeticRepository.FindByID(ctx, cosmeticID)
	if err != nil {
		return nil, err
	}
	if cosmeticEntity == nil {
		return nil, domainerror.New("COSMETIC_NOT_FOUND", "cosmético no encontrado", domainerror.ErrNotFound)
	}
	return cosmeticEntity, nil
}
