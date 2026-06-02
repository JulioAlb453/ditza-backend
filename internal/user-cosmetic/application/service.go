package application

import (
	"context"

	cosmeticdomain "ditza/internal/cosmetic/domain"
	pointtransactiondomain "ditza/internal/point-transaction/domain"
	domainerror "ditza/internal/shared/domain/error"
	unitofwork "ditza/internal/shared/domain/unit-of-work"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	usercosmeticdomain "ditza/internal/user-cosmetic/domain"
	userdomain "ditza/internal/user/domain"
)

type Service struct {
	unitOfWork                 unitofwork.UnitOfWork
	userRepository             userdomain.Repository
	cosmeticRepository         cosmeticdomain.Repository
	userCosmeticRepository     usercosmeticdomain.Repository
	pointTransactionRepository pointtransactiondomain.Repository
}

type BuyCosmeticCommand struct {
	UserID     valueobject.UserID
	CosmeticID valueobject.CosmeticID
}

type BuyCosmeticResult struct {
	WalletCoins int
	CosmeticID  valueobject.CosmeticID
}

func NewService(
	unitOfWork unitofwork.UnitOfWork,
	userRepository userdomain.Repository,
	cosmeticRepository cosmeticdomain.Repository,
	userCosmeticRepository usercosmeticdomain.Repository,
	pointTransactionRepository pointtransactiondomain.Repository,
) *Service {
	return &Service{
		unitOfWork:                 unitOfWork,
		userRepository:             userRepository,
		cosmeticRepository:         cosmeticRepository,
		userCosmeticRepository:     userCosmeticRepository,
		pointTransactionRepository: pointTransactionRepository,
	}
}

func (s *Service) Buy(ctx context.Context, command BuyCosmeticCommand) (result *BuyCosmeticResult, err error) {
	tracker := monitoring.StartService(logger.ModelUserCosmetic, "buy", map[string]any{
		"user_id":     command.UserID,
		"cosmetic_id": command.CosmeticID,
	})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"wallet_coins": result.WalletCoins})
	}()

	result = &BuyCosmeticResult{}
	err = s.withinTransaction(ctx, func(txCtx context.Context) error {
		userEntity, err := s.userRepository.FindByID(txCtx, command.UserID)
		if err != nil {
			return err
		}
		if userEntity == nil {
			return domainerror.New("USER_NOT_FOUND", "usuario no encontrado", domainerror.ErrNotFound)
		}

		cosmeticEntity, err := s.cosmeticRepository.FindByID(txCtx, command.CosmeticID)
		if err != nil {
			return err
		}
		if cosmeticEntity == nil {
			return domainerror.New("COSMETIC_NOT_FOUND", "cosmético no encontrado", domainerror.ErrNotFound)
		}
		if !cosmeticEntity.CanBePurchased() {
			return domainerror.New("COSMETIC_NOT_AVAILABLE", "el cosmético no está disponible", domainerror.ErrInvalidInput)
		}

		alreadyOwned, err := s.userCosmeticRepository.Exists(txCtx, command.UserID, command.CosmeticID)
		if err != nil {
			return err
		}
		if alreadyOwned {
			return domainerror.New("COSMETIC_ALREADY_OWNED", "ya posees este cosmético", domainerror.ErrCosmeticAlreadyOwned)
		}

		if err := userEntity.SpendCoins(cosmeticEntity.PriceCoins); err != nil {
			return err
		}
		if err := s.userRepository.Update(txCtx, userEntity); err != nil {
			return err
		}

		userCosmetic := usercosmeticdomain.New(command.UserID, command.CosmeticID)
		if err := s.userCosmeticRepository.Create(txCtx, userCosmetic); err != nil {
			return err
		}

		pointTransaction, err := pointtransactiondomain.New(
			command.UserID,
			pointtransactiondomain.TypePurchase,
			-cosmeticEntity.PriceCoins,
			0,
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.pointTransactionRepository.Create(txCtx, pointTransaction); err != nil {
			return err
		}

		result.WalletCoins = userEntity.Coins
		result.CosmeticID = command.CosmeticID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) ListInventory(ctx context.Context, userID valueobject.UserID) (inventory []usercosmeticdomain.UserCosmetic, err error) {
	tracker := monitoring.StartService(logger.ModelUserCosmetic, "list_inventory", map[string]any{"user_id": userID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"count": len(inventory)})
	}()

	inventory, err = s.userCosmeticRepository.FindByUserID(ctx, userID)
	return inventory, err
}

func (s *Service) withinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if s.unitOfWork == nil {
		return fn(ctx)
	}
	return s.unitOfWork.WithinTransaction(ctx, fn)
}
