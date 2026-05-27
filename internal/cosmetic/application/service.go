package application

import (
	"context"

	cosmeticdomain "ditza/internal/cosmetic/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Service struct {
	cosmeticRepository cosmeticdomain.Repository
}

func NewService(cosmeticRepository cosmeticdomain.Repository) *Service {
	return &Service{cosmeticRepository: cosmeticRepository}
}

func (s *Service) ListActive(ctx context.Context) ([]cosmeticdomain.Cosmetic, error) {
	return s.cosmeticRepository.FindAllActive(ctx)
}

func (s *Service) GetByID(ctx context.Context, cosmeticID valueobject.CosmeticID) (*cosmeticdomain.Cosmetic, error) {
	cosmeticEntity, err := s.cosmeticRepository.FindByID(ctx, cosmeticID)
	if err != nil {
		return nil, err
	}
	if cosmeticEntity == nil {
		return nil, domainerror.New("COSMETIC_NOT_FOUND", "cosmético no encontrado", domainerror.ErrNotFound)
	}
	return cosmeticEntity, nil
}
