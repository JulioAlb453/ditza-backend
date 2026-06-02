package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ditza/internal/season/data"
	seasondomain "ditza/internal/season/domain"
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

func (r *Repository) Create(ctx context.Context, entity *seasondomain.Season) error {
	monitoring.Repository(logger.ModelSeason, "create", map[string]any{"starts_at": entity.StartsAt})

	model := data.ToModel(*entity)
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		INSERT INTO seasons (starts_at, ends_at, is_active, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		model.StartsAt,
		model.EndsAt,
		model.IsActive,
		model.CreatedAt,
	).Scan(&model.ID)
	if err != nil {
		return fmt.Errorf("error creando temporada: %w", err)
	}

	entity.ID = valueobject.SeasonID(model.ID)
	return nil
}

func (r *Repository) Update(ctx context.Context, entity *seasondomain.Season) error {
	monitoring.Repository(logger.ModelSeason, "update", map[string]any{"season_id": entity.ID})

	model := data.ToModel(*entity)
	result, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE seasons
		SET starts_at = $2, ends_at = $3, is_active = $4
		WHERE id = $1`,
		model.ID,
		model.StartsAt,
		model.EndsAt,
		model.IsActive,
	)
	if err != nil {
		return fmt.Errorf("error actualizando temporada: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo verificar la actualización de temporada: %w", err)
	}
	if rowsAffected == 0 {
		return domainerror.New("SEASON_NOT_FOUND", "temporada no encontrada", domainerror.ErrNotFound)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id valueobject.SeasonID) (*seasondomain.Season, error) {
	monitoring.Repository(logger.ModelSeason, "find_by_id", map[string]any{"season_id": id})

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

func (r *Repository) FindActive(ctx context.Context) (*seasondomain.Season, error) {
	monitoring.Repository(logger.ModelSeason, "find_active", nil)

	model, err := r.queryOne(ctx, "WHERE is_active = true ORDER BY created_at DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}

	entity := data.ToDomain(*model)
	return &entity, nil
}

func (r *Repository) DeactivateAll(ctx context.Context) error {
	monitoring.Repository(logger.ModelSeason, "deactivate_all", nil)

	_, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE seasons SET is_active = false WHERE is_active = true`)
	if err != nil {
		return fmt.Errorf("error desactivando temporadas: %w", err)
	}
	return nil
}

func (r *Repository) queryOne(ctx context.Context, whereClause string, args ...any) (*data.Model, error) {
	query := fmt.Sprintf(`
		SELECT id, starts_at, ends_at, is_active, created_at
		FROM seasons %s`, whereClause)

	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, query, args...).Scan(
		&model.ID,
		&model.StartsAt,
		&model.EndsAt,
		&model.IsActive,
		&model.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando temporada: %w", err)
	}
	return &model, nil
}
