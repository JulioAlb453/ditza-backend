package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"ditza/internal/season-score/data"
	seasonscoredomain "ditza/internal/season-score/domain"
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

func (r *Repository) Create(ctx context.Context, entity *seasonscoredomain.SeasonScore) error {
	monitoring.Repository(logger.ModelSeasonScore, "create", map[string]any{
		"user_id":   entity.UserID,
		"season_id": entity.SeasonID,
	})

	model := data.ToModel(*entity)
	_, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		INSERT INTO season_scores (user_id, season_id, points, updated_at)
		VALUES ($1, $2, $3, $4)`,
		model.UserID,
		model.SeasonID,
		model.Points,
		model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creando puntaje de temporada: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, entity *seasonscoredomain.SeasonScore) error {
	monitoring.Repository(logger.ModelSeasonScore, "update", map[string]any{
		"user_id":   entity.UserID,
		"season_id": entity.SeasonID,
	})

	model := data.ToModel(*entity)
	result, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE season_scores
		SET points = $3, updated_at = $4
		WHERE user_id = $1 AND season_id = $2`,
		model.UserID,
		model.SeasonID,
		model.Points,
		model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error actualizando puntaje de temporada: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo verificar la actualización del puntaje: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("puntaje de temporada no encontrado")
	}
	return nil
}

func (r *Repository) FindByUserAndSeason(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) (*seasonscoredomain.SeasonScore, error) {
	monitoring.Repository(logger.ModelSeasonScore, "find_by_user_and_season", map[string]any{
		"user_id":   userID,
		"season_id": seasonID,
	})

	model, err := r.queryOne(ctx, userID.String(), uint64(seasonID))
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}

	entity := data.ToDomain(*model)
	return &entity, nil
}

func (r *Repository) FindByUserIDsAndSeason(ctx context.Context, userIDs []valueobject.UserID, seasonID valueobject.SeasonID) ([]seasonscoredomain.SeasonScore, error) {
	monitoring.Repository(logger.ModelSeasonScore, "find_by_user_ids_and_season", map[string]any{
		"user_count": len(userIDs),
		"season_id":  seasonID,
	})

	if len(userIDs) == 0 {
		return []seasonscoredomain.SeasonScore{}, nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]any, 0, len(userIDs)+1)
	for i, userID := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args = append(args, userID.String())
	}
	args = append(args, uint64(seasonID))

	query := fmt.Sprintf(`
		SELECT user_id, season_id, points, updated_at
		FROM season_scores
		WHERE user_id IN (%s) AND season_id = $%d`,
		strings.Join(placeholders, ", "), len(args),
	)

	rows, err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listando puntajes de temporada: %w", err)
	}
	defer rows.Close()

	var models []data.Model
	for rows.Next() {
		var model data.Model
		if err := rows.Scan(&model.UserID, &model.SeasonID, &model.Points, &model.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error leyendo puntaje de temporada: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando puntajes de temporada: %w", err)
	}
	return data.ToDomainList(models), nil
}

func (r *Repository) ResetBySeasonID(ctx context.Context, seasonID valueobject.SeasonID) error {
	monitoring.Repository(logger.ModelSeasonScore, "reset_by_season_id", map[string]any{"season_id": seasonID})

	_, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE season_scores
		SET points = 0, updated_at = NOW()
		WHERE season_id = $1`, uint64(seasonID))
	if err != nil {
		return fmt.Errorf("error reiniciando puntajes de temporada: %w", err)
	}
	return nil
}

func (r *Repository) queryOne(ctx context.Context, userID string, seasonID uint64) (*data.Model, error) {
	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT user_id, season_id, points, updated_at
		FROM season_scores
		WHERE user_id = $1 AND season_id = $2`, userID, seasonID,
	).Scan(&model.UserID, &model.SeasonID, &model.Points, &model.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando puntaje de temporada: %w", err)
	}
	return &model, nil
}
