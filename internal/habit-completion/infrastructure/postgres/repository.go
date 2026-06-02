package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ditza/internal/habit-completion/data"
	habitcompletiondomain "ditza/internal/habit-completion/domain"
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

func (r *Repository) Create(ctx context.Context, entity *habitcompletiondomain.HabitCompletion) error {
	monitoring.Repository(logger.ModelHabitCompletion, "create", map[string]any{
		"habit_id": entity.HabitID,
		"user_id":  entity.UserID,
	})

	model := data.ToModel(*entity)
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		INSERT INTO habit_completions (
			habit_id, user_id, completed_at, note, emoji,
			coins_awarded, season_points_awarded, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`,
		model.HabitID,
		model.UserID,
		model.CompletedAt,
		model.Note,
		model.Emoji,
		model.CoinsAwarded,
		model.SeasonPointsAwarded,
		model.CreatedAt,
	).Scan(&model.ID)
	if err != nil {
		return fmt.Errorf("error creando completado de hábito: %w", err)
	}

	entity.ID = valueobject.HabitCompletionID(model.ID)
	return nil
}

func (r *Repository) ExistsForHabitOnDate(ctx context.Context, habitID valueobject.HabitID, date time.Time) (bool, error) {
	monitoring.Repository(logger.ModelHabitCompletion, "exists_for_habit_on_date", map[string]any{
		"habit_id": habitID,
		"date":     date.Format("2006-01-02"),
	})

	var exists bool
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM habit_completions
			WHERE habit_id = $1
			  AND completion_date = ($2 AT TIME ZONE 'UTC')::date
		)`, uint64(habitID), date,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando completado diario: %w", err)
	}
	return exists, nil
}

func (r *Repository) CountByUserOnDate(ctx context.Context, userID valueobject.UserID, date time.Time) (int, error) {
	monitoring.Repository(logger.ModelHabitCompletion, "count_by_user_on_date", map[string]any{
		"user_id": userID,
		"date":    date.Format("2006-01-02"),
	})

	var count int
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT COUNT(*) FROM habit_completions
		WHERE user_id = $1
		  AND completion_date = ($2 AT TIME ZONE 'UTC')::date`,
		userID.String(), date,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error contando completados del usuario: %w", err)
	}
	return count, nil
}

func (r *Repository) FindByUserIDAndDateRange(ctx context.Context, userID valueobject.UserID, from, to time.Time) ([]habitcompletiondomain.HabitCompletion, error) {
	monitoring.Repository(logger.ModelHabitCompletion, "find_by_user_id_and_date_range", map[string]any{
		"user_id": userID,
		"from":    from.Format("2006-01-02"),
		"to":      to.Format("2006-01-02"),
	})

	rows, err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryContext(ctx, `
		SELECT id, habit_id, user_id, completed_at, note, emoji,
		       coins_awarded, season_points_awarded, created_at
		FROM habit_completions
		WHERE user_id = $1
		  AND completion_date BETWEEN ($2 AT TIME ZONE 'UTC')::date AND ($3 AT TIME ZONE 'UTC')::date
		ORDER BY completed_at ASC`,
		userID.String(), from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("error listando completados: %w", err)
	}
	defer rows.Close()

	var models []data.Model
	for rows.Next() {
		var model data.Model
		if err := rows.Scan(
			&model.ID,
			&model.HabitID,
			&model.UserID,
			&model.CompletedAt,
			&model.Note,
			&model.Emoji,
			&model.CoinsAwarded,
			&model.SeasonPointsAwarded,
			&model.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error leyendo completado: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando completados: %w", err)
	}
	return data.ToDomainList(models), nil
}
