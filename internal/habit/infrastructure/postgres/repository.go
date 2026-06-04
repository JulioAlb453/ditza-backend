package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ditza/internal/habit/data"
	habitdomain "ditza/internal/habit/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	sharedpostgres "ditza/internal/shared/infrastructure/postgres"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, entity *habitdomain.Habit) error {
	monitoring.Repository(logger.ModelHabit, "create", map[string]any{
		"user_id": entity.UserID,
		"title":   entity.Title,
	})

	model := data.ToModel(*entity)
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		INSERT INTO habits (
			user_id, title, description, emoji, category, color, frequency,
			target_count, target_unit, difficulty, reminder_time, is_active,
			current_streak, best_streak, last_completed_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id`,
		model.UserID,
		model.Title,
		model.Description,
		model.Emoji,
		model.Category,
		model.Color,
		model.Frequency,
		model.TargetCount,
		model.TargetUnit,
		model.Difficulty,
		model.ReminderTime,
		model.IsActive,
		model.CurrentStreak,
		model.BestStreak,
		model.LastCompletedDate,
		model.CreatedAt,
		model.UpdatedAt,
	).Scan(&model.ID)
	if err != nil {
		return fmt.Errorf("error creando hábito: %w", err)
	}

	entity.ID = valueobject.HabitID(model.ID)
	return nil
}

func (r *Repository) Update(ctx context.Context, entity *habitdomain.Habit) error {
	monitoring.Repository(logger.ModelHabit, "update", map[string]any{"habit_id": entity.ID})

	model := data.ToModel(*entity)
	result, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE habits
		SET title = $2, description = $3, emoji = $4, category = $5, color = $6,
		    frequency = $7, target_count = $8, target_unit = $9, difficulty = $10,
		    reminder_time = $11, is_active = $12, current_streak = $13, best_streak = $14,
		    last_completed_date = $15, updated_at = $16
		WHERE id = $1`,
		model.ID,
		model.Title,
		model.Description,
		model.Emoji,
		model.Category,
		model.Color,
		model.Frequency,
		model.TargetCount,
		model.TargetUnit,
		model.Difficulty,
		model.ReminderTime,
		model.IsActive,
		model.CurrentStreak,
		model.BestStreak,
		model.LastCompletedDate,
		model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error actualizando hábito: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo verificar la actualización del hábito: %w", err)
	}
	if rowsAffected == 0 {
		return domainerror.New("HABIT_NOT_FOUND", "hábito no encontrado", domainerror.ErrNotFound)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id valueobject.HabitID) (*habitdomain.Habit, error) {
	monitoring.Repository(logger.ModelHabit, "find_by_id", map[string]any{"habit_id": id})

	model, err := r.queryOne(ctx, "WHERE id = $1", uint64(id))
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}

	entity := data.ToDomain(*model)
	return &entity, nil
}

func (r *Repository) FindActiveByUserID(ctx context.Context, userID valueobject.UserID) ([]habitdomain.Habit, error) {
	monitoring.Repository(logger.ModelHabit, "find_active_by_user_id", map[string]any{"user_id": userID})

	return r.queryList(ctx, `
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at ASC`, userID.String())
}

func (r *Repository) CountActiveByUserID(ctx context.Context, userID valueobject.UserID) (int, error) {
	monitoring.Repository(logger.ModelHabit, "count_active_by_user_id", map[string]any{"user_id": userID})

	var count int
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT COUNT(*) FROM habits WHERE user_id = $1 AND is_active = true`, userID.String(),
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error contando hábitos activos: %w", err)
	}
	return count, nil
}

func (r *Repository) queryOne(ctx context.Context, whereClause string, arg any) (*data.Model, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, title, description, emoji, category, color, frequency,
		       target_count, target_unit, difficulty, reminder_time, is_active,
		       current_streak, best_streak, last_completed_date, created_at, updated_at
		FROM habits %s`, whereClause)

	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, query, arg).Scan(
		&model.ID,
		&model.UserID,
		&model.Title,
		&model.Description,
		&model.Emoji,
		&model.Category,
		&model.Color,
		&model.Frequency,
		&model.TargetCount,
		&model.TargetUnit,
		&model.Difficulty,
		&model.ReminderTime,
		&model.IsActive,
		&model.CurrentStreak,
		&model.BestStreak,
		&model.LastCompletedDate,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando hábito: %w", err)
	}
	return &model, nil
}

func (r *Repository) queryList(ctx context.Context, whereClause string, args ...any) ([]habitdomain.Habit, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, title, description, emoji, category, color, frequency,
		       target_count, target_unit, difficulty, reminder_time, is_active,
		       current_streak, best_streak, last_completed_date, created_at, updated_at
		FROM habits %s`, whereClause)

	rows, err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listando hábitos: %w", err)
	}
	defer rows.Close()

	var models []data.Model
	for rows.Next() {
		var model data.Model
		if err := rows.Scan(
			&model.ID,
			&model.UserID,
			&model.Title,
			&model.Description,
			&model.Emoji,
			&model.Category,
			&model.Color,
			&model.Frequency,
			&model.TargetCount,
			&model.TargetUnit,
			&model.Difficulty,
			&model.ReminderTime,
			&model.IsActive,
			&model.CurrentStreak,
			&model.BestStreak,
			&model.LastCompletedDate,
			&model.CreatedAt,
			&model.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error leyendo hábito: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando hábitos: %w", err)
	}
	return data.ToDomainList(models), nil
}
