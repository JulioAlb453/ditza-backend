package application

import (
	"context"
	"time"

	habitcompletiondomain "ditza/internal/habit-completion/domain"
	habitdomain "ditza/internal/habit/domain"
	petdomain "ditza/internal/pet/domain"
	pointtransactiondomain "ditza/internal/point-transaction/domain"
	seasonscoredomain "ditza/internal/season-score/domain"
	seasondomain "ditza/internal/season/domain"
	domainerror "ditza/internal/shared/domain/error"
	unitofwork "ditza/internal/shared/domain/unit-of-work"
	valueobject "ditza/internal/shared/domain/value-object"
	userdomain "ditza/internal/user/domain"
)

type Service struct {
	unitOfWork                 unitofwork.UnitOfWork
	habitRepository            habitdomain.Repository
	habitCompletionRepository  habitcompletiondomain.Repository
	userRepository             userdomain.Repository
	seasonRepository           seasondomain.Repository
	seasonScoreRepository      seasonscoredomain.Repository
	petRepository              petdomain.Repository
	pointTransactionRepository pointtransactiondomain.Repository
}

type CompleteHabitCommand struct {
	UserID      valueobject.UserID
	HabitID     valueobject.HabitID
	CompletedAt time.Time
	Note        *string
	Emoji       *string
}

type CompleteHabitResult struct {
	CoinsEarned         int
	SeasonPointsEarned  int
	CurrentStreak       int
	WalletCoins         int
	CurrentSeasonPoints int
	PetLevel            int
	PetMood             string
}

func NewService(
	unitOfWork unitofwork.UnitOfWork,
	habitRepository habitdomain.Repository,
	habitCompletionRepository habitcompletiondomain.Repository,
	userRepository userdomain.Repository,
	seasonRepository seasondomain.Repository,
	seasonScoreRepository seasonscoredomain.Repository,
	petRepository petdomain.Repository,
	pointTransactionRepository pointtransactiondomain.Repository,
) *Service {
	return &Service{
		unitOfWork:                 unitOfWork,
		habitRepository:            habitRepository,
		habitCompletionRepository:  habitCompletionRepository,
		userRepository:             userRepository,
		seasonRepository:           seasonRepository,
		seasonScoreRepository:      seasonScoreRepository,
		petRepository:              petRepository,
		pointTransactionRepository: pointTransactionRepository,
	}
}

func (s *Service) Complete(ctx context.Context, command CompleteHabitCommand) (*CompleteHabitResult, error) {
	completedAt := command.CompletedAt.UTC()
	if completedAt.IsZero() {
		completedAt = time.Now().UTC()
	}

	result := &CompleteHabitResult{}
	err := s.withinTransaction(ctx, func(txCtx context.Context) error {
		habitEntity, err := s.habitRepository.FindByID(txCtx, command.HabitID)
		if err != nil {
			return err
		}
		if habitEntity == nil {
			return domainerror.New("HABIT_NOT_FOUND", "hábito no encontrado", domainerror.ErrNotFound)
		}
		if !habitEntity.BelongsTo(command.UserID) {
			return domainerror.New("HABIT_NOT_OWNED", "el hábito no pertenece al usuario", domainerror.ErrHabitNotOwned)
		}
		if err := habitEntity.CanCompleteOn(completedAt); err != nil {
			return err
		}

		alreadyCompleted, err := s.habitCompletionRepository.ExistsForHabitOnDate(txCtx, command.HabitID, completedAt)
		if err != nil {
			return err
		}
		if alreadyCompleted {
			return domainerror.New("HABIT_ALREADY_COMPLETED", "el hábito ya fue completado hoy", domainerror.ErrHabitAlreadyCompleted)
		}

		habitEntity.ApplyCompletion(completedAt)
		if err := s.habitRepository.Update(txCtx, habitEntity); err != nil {
			return err
		}

		hasJournalEntry := command.Note != nil || command.Emoji != nil
		reward := habitcompletiondomain.CalculateCompletionReward(habitEntity.CurrentStreak, hasJournalEntry)
		completionEntity, err := habitcompletiondomain.New(
			command.HabitID,
			command.UserID,
			completedAt,
			command.Note,
			command.Emoji,
			reward,
		)
		if err != nil {
			return err
		}
		if err := s.habitCompletionRepository.Create(txCtx, completionEntity); err != nil {
			return err
		}

		userEntity, err := s.userRepository.FindByID(txCtx, command.UserID)
		if err != nil {
			return err
		}
		if userEntity == nil {
			return domainerror.New("USER_NOT_FOUND", "usuario no encontrado", domainerror.ErrNotFound)
		}
		if err := userEntity.AddCoins(reward.Coins); err != nil {
			return err
		}
		if err := s.userRepository.Update(txCtx, userEntity); err != nil {
			return err
		}

		seasonEntity, err := s.seasonRepository.FindActive(txCtx)
		if err != nil {
			return err
		}
		if seasonEntity == nil {
			return domainerror.New("SEASON_NOT_ACTIVE", "no hay una temporada activa", domainerror.ErrSeasonNotActive)
		}

		seasonScoreEntity, err := s.seasonScoreRepository.FindByUserAndSeason(txCtx, command.UserID, seasonEntity.ID)
		if err != nil {
			return err
		}
		if seasonScoreEntity == nil {
			seasonScoreEntity = seasonscoredomain.New(command.UserID, seasonEntity.ID)
			if err := s.seasonScoreRepository.Create(txCtx, seasonScoreEntity); err != nil {
				return err
			}
		}
		if err := seasonScoreEntity.AddPoints(reward.SeasonPoints); err != nil {
			return err
		}
		if err := s.seasonScoreRepository.Update(txCtx, seasonScoreEntity); err != nil {
			return err
		}

		petEntity, err := s.petRepository.FindByUserID(txCtx, command.UserID)
		if err != nil {
			return err
		}
		if petEntity == nil {
			petEntity, err = petdomain.New(command.UserID, "Ditza")
			if err != nil {
				return err
			}
			if err := s.petRepository.Create(txCtx, petEntity); err != nil {
				return err
			}
		}

		activeHabits, err := s.habitRepository.FindActiveByUserID(txCtx, command.UserID)
		if err != nil {
			return err
		}
		completedToday, err := s.habitCompletionRepository.CountByUserOnDate(txCtx, command.UserID, completedAt)
		if err != nil {
			return err
		}

		petEntity.AddXP(reward.SeasonPoints)
		petEntity.RegisterInteraction(completedAt)
		petEntity.UpdateMoodFromProgress(completedToday, len(activeHabits))
		if err := s.petRepository.Update(txCtx, petEntity); err != nil {
			return err
		}

		txType := pointtransactiondomain.TypeHabit
		if reward.NoteBonus > 0 {
			txType = pointtransactiondomain.TypeNoteBonus
		} else if reward.StreakBonus > 0 {
			txType = pointtransactiondomain.TypeStreakBonus
		}
		pointTransaction, err := pointtransactiondomain.New(
			command.UserID,
			txType,
			reward.Coins,
			reward.SeasonPoints,
			nil,
		)
		if err != nil {
			return err
		}
		if err := s.pointTransactionRepository.Create(txCtx, pointTransaction); err != nil {
			return err
		}

		result.CoinsEarned = reward.Coins
		result.SeasonPointsEarned = reward.SeasonPoints
		result.CurrentStreak = habitEntity.CurrentStreak
		result.WalletCoins = userEntity.Coins
		result.CurrentSeasonPoints = seasonScoreEntity.Points
		result.PetLevel = petEntity.Level
		result.PetMood = petEntity.Mood.String()
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) withinTransaction(ctx context.Context, fn func(context.Context) error) error {
	if s.unitOfWork == nil {
		return fn(ctx)
	}
	return s.unitOfWork.WithinTransaction(ctx, fn)
}
