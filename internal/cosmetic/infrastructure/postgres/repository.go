package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ditza/internal/cosmetic/data"
	cosmeticdomain "ditza/internal/cosmetic/domain"
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

func (r *Repository) Create(ctx context.Context, entity *cosmeticdomain.Cosmetic) error {
	monitoring.Repository(logger.ModelCosmetic, "create", map[string]any{"name": entity.Name})

	model := data.ToModel(*entity)
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		INSERT INTO cosmetics (name, slot, price_coins, rarity, asset_key, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		model.Name,
		model.Slot,
		model.PriceCoins,
		model.Rarity,
		model.AssetKey,
		model.IsActive,
		model.CreatedAt,
	).Scan(&model.ID)
	if err != nil {
		return fmt.Errorf("error creando cosmético: %w", err)
	}

	entity.ID = valueobject.CosmeticID(model.ID)
	return nil
}

func (r *Repository) Update(ctx context.Context, entity *cosmeticdomain.Cosmetic) error {
	monitoring.Repository(logger.ModelCosmetic, "update", map[string]any{"cosmetic_id": entity.ID})

	model := data.ToModel(*entity)
	result, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE cosmetics
		SET name = $2, slot = $3, price_coins = $4, rarity = $5, asset_key = $6, is_active = $7
		WHERE id = $1`,
		model.ID,
		model.Name,
		model.Slot,
		model.PriceCoins,
		model.Rarity,
		model.AssetKey,
		model.IsActive,
	)
	if err != nil {
		return fmt.Errorf("error actualizando cosmético: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo verificar la actualización del cosmético: %w", err)
	}
	if rowsAffected == 0 {
		return domainerror.New("COSMETIC_NOT_FOUND", "cosmético no encontrado", domainerror.ErrNotFound)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id valueobject.CosmeticID) (*cosmeticdomain.Cosmetic, error) {
	monitoring.Repository(logger.ModelCosmetic, "find_by_id", map[string]any{"cosmetic_id": id})

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

func (r *Repository) FindAllActive(ctx context.Context) ([]cosmeticdomain.Cosmetic, error) {
	monitoring.Repository(logger.ModelCosmetic, "find_all_active", nil)

	rows, err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryContext(ctx, `
		SELECT id, name, slot, price_coins, rarity, asset_key, is_active, created_at
		FROM cosmetics
		WHERE is_active = true
		ORDER BY price_coins ASC, name ASC`)
	if err != nil {
		return nil, fmt.Errorf("error listando cosméticos: %w", err)
	}
	defer rows.Close()

	var models []data.Model
	for rows.Next() {
		var model data.Model
		if err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.Slot,
			&model.PriceCoins,
			&model.Rarity,
			&model.AssetKey,
			&model.IsActive,
			&model.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error leyendo cosmético: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando cosméticos: %w", err)
	}
	return data.ToDomainList(models), nil
}

func (r *Repository) queryOne(ctx context.Context, whereClause string, arg any) (*data.Model, error) {
	query := fmt.Sprintf(`
		SELECT id, name, slot, price_coins, rarity, asset_key, is_active, created_at
		FROM cosmetics %s`, whereClause)

	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, query, arg).Scan(
		&model.ID,
		&model.Name,
		&model.Slot,
		&model.PriceCoins,
		&model.Rarity,
		&model.AssetKey,
		&model.IsActive,
		&model.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando cosmético: %w", err)
	}
	return &model, nil
}
