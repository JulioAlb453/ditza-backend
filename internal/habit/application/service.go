package application

import (
	"context"

	habitdomain "ditza/internal/habit/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

type Service struct {
	habitRepository habitdomain.Repository
}

type CreateHabitCommand struct {
	UserID       valueobject.UserID
	Title        string
	Description  string
	Emoji        string
	Category     string
	Color        string
	Frequency    string
	TargetCount  int
	TargetUnit   string
	Difficulty   string
	ReminderTime *string
}

func NewService(habitRepository habitdomain.Repository) *Service {
	return &Service{habitRepository: habitRepository}
}

func (s *Service) Create(ctx context.Context, command CreateHabitCommand) (habitEntity *habitdomain.Habit, err error) {
	tracker := monitoring.StartService(logger.ModelHabit, "create", map[string]any{"user_id": command.UserID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"habit_id": habitEntity.ID})
	}()

	activeCount, err := s.habitRepository.CountActiveByUserID(ctx, command.UserID)
	if err != nil {
		return nil, err
	}
	if activeCount >= valueobject.MaxActiveHabitsPerUser {
		return nil, domainerror.New("HABIT_LIMIT_REACHED", "límite de hábitos activos alcanzado", domainerror.ErrHabitLimitReached)
	}

	habitEntity, err = habitdomain.New(command.UserID, habitdomain.HabitConfig{
		Title:        command.Title,
		Description:  command.Description,
		Emoji:        command.Emoji,
		Category:     command.Category,
		Color:        command.Color,
		Frequency:    command.Frequency,
		TargetCount:  command.TargetCount,
		TargetUnit:   command.TargetUnit,
		Difficulty:   command.Difficulty,
		ReminderTime: command.ReminderTime,
	})
	if err != nil {
		return nil, err
	}

	if err := s.habitRepository.Create(ctx, habitEntity); err != nil {
		return nil, err
	}

	return habitEntity, nil
}

func (s *Service) ListActiveByUser(ctx context.Context, userID valueobject.UserID) (habits []habitdomain.Habit, err error) {
	tracker := monitoring.StartService(logger.ModelHabit, "list_active_by_user", map[string]any{"user_id": userID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(map[string]any{"count": len(habits)})
	}()

	habits, err = s.habitRepository.FindActiveByUserID(ctx, userID)
	return habits, err
}

func (s *Service) Deactivate(ctx context.Context, userID valueobject.UserID, habitID valueobject.HabitID) (err error) {
	tracker := monitoring.StartService(logger.ModelHabit, "deactivate", map[string]any{"user_id": userID, "habit_id": habitID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(nil)
	}()

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
