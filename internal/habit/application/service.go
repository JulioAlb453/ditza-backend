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
	Category     string
	Color        string
	Frequency    string
	TargetCount  int
	TargetUnit   string
	Difficulty   string
	ReminderTime string
}

type UpdateHabitCommand struct {
	UserID       valueobject.UserID
	HabitID      valueobject.HabitID
	Title        *string
	Description  *string
	Category     *string
	Color        *string
	Frequency    *string
	TargetCount  *int
	TargetUnit   *string
	Difficulty   *string
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

	habitEntity, err = habitdomain.New(command.UserID, habitdomain.HabitConfig{
		Title:        command.Title,
		Description:  command.Description,
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

func (s *Service) Update(ctx context.Context, command UpdateHabitCommand) (habitEntity *habitdomain.Habit, err error) {
	tracker := monitoring.StartService(logger.ModelHabit, "update", map[string]any{"user_id": command.UserID, "habit_id": command.HabitID})
	defer func() {
		if err != nil {
			tracker.Fail(err, nil)
			return
		}
		tracker.Success(nil)
	}()

	habitEntity, err = s.habitRepository.FindByID(ctx, command.HabitID)
	if err != nil {
		return nil, err
	}
	if habitEntity == nil {
		return nil, domainerror.New("HABIT_NOT_FOUND", "hábito no encontrado", domainerror.ErrNotFound)
	}
	if !habitEntity.BelongsTo(command.UserID) {
		return nil, domainerror.New("HABIT_NOT_OWNED", "el hábito no pertenece al usuario", domainerror.ErrHabitNotOwned)
	}

	config := habitdomain.HabitConfig{
		Title:        habitEntity.Title,
		Description:  habitEntity.Description,
		Category:     habitEntity.Category,
		Color:        habitEntity.Color,
		Frequency:    habitEntity.Frequency,
		TargetCount:  habitEntity.TargetCount,
		TargetUnit:   habitEntity.TargetUnit,
		Difficulty:   habitEntity.Difficulty,
		ReminderTime: habitEntity.ReminderTime,
	}
	if command.Title != nil {
		config.Title = *command.Title
	}
	if command.Description != nil {
		config.Description = *command.Description
	}
	if command.Category != nil {
		config.Category = *command.Category
	}
	if command.Color != nil {
		config.Color = *command.Color
	}
	if command.Frequency != nil {
		config.Frequency = *command.Frequency
	}
	if command.TargetCount != nil {
		config.TargetCount = *command.TargetCount
	}
	if command.TargetUnit != nil {
		config.TargetUnit = *command.TargetUnit
	}
	if command.Difficulty != nil {
		config.Difficulty = *command.Difficulty
	}
	if command.ReminderTime != nil {
		config.ReminderTime = *command.ReminderTime
	}

	if err := habitEntity.Update(config); err != nil {
		return nil, err
	}
	if err := s.habitRepository.Update(ctx, habitEntity); err != nil {
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
