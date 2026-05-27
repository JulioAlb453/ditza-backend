package application

import (
	"context"

	habitdomain "ditza/internal/habit/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
)

type Service struct {
	habitRepository habitdomain.Repository
}

type CreateHabitCommand struct {
	UserID valueobject.UserID
	Title  string
}

func NewService(habitRepository habitdomain.Repository) *Service {
	return &Service{habitRepository: habitRepository}
}

func (s *Service) Create(ctx context.Context, command CreateHabitCommand) (*habitdomain.Habit, error) {
	activeCount, err := s.habitRepository.CountActiveByUserID(ctx, command.UserID)
	if err != nil {
		return nil, err
	}
	if activeCount >= valueobject.MaxActiveHabitsPerUser {
		return nil, domainerror.New("HABIT_LIMIT_REACHED", "límite de hábitos activos alcanzado", domainerror.ErrHabitLimitReached)
	}

	habitEntity, err := habitdomain.New(command.UserID, command.Title)
	if err != nil {
		return nil, err
	}

	if err := s.habitRepository.Create(ctx, habitEntity); err != nil {
		return nil, err
	}

	return habitEntity, nil
}

func (s *Service) ListActiveByUser(ctx context.Context, userID valueobject.UserID) ([]habitdomain.Habit, error) {
	return s.habitRepository.FindActiveByUserID(ctx, userID)
}

func (s *Service) Deactivate(ctx context.Context, userID valueobject.UserID, habitID valueobject.HabitID) error {
	habitEntity, err := s.habitRepository.FindByID(ctx, habitID)
	if err != nil {
		return err
	}
	if habitEntity == nil {
		return domainerror.New("HABIT_NOT_FOUND", "hábito no encontrado", domainerror.ErrNotFound)
	}
	if !habitEntity.BelongsTo(userID) {
		return domainerror.New("HABIT_NOT_OWNED", "el hábito no pertenece al usuario", domainerror.ErrHabitNotOwned)
	}

	habitEntity.Deactivate()
	return s.habitRepository.Update(ctx, habitEntity)
}
