package postgres

import (
	"context"
	"database/sql"
	"fmt"

	usercosmeticdata "ditza/internal/user-cosmetic/data"
	usercosmeticdomain "ditza/internal/user-cosmetic/domain"
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

func (r *Repository) Create(ctx context.Context, entity *usercosmeticdomain.UserCosmetic) error {
	monitoring.Repository(logger.ModelUserCosmetic, "create", map[string]any{
		"user_id":     entity.UserID,
		"cosmetic_id": entity.CosmeticID,
	})

	model := usercosmeticdata.ToModel(*entity)
	_, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		INSERT INTO user_cosmetics (user_id, cosmetic_id, purchased_at)
		VALUES ($1, $2, $3)`,
		model.UserID,
		model.CosmeticID,
		model.PurchasedAt,
	)
	if err != nil {
		return fmt.Errorf("error registrando cosmético del usuario: %w", err)
	}
	return nil
}

func (r *Repository) Exists(ctx context.Context, userID valueobject.UserID, cosmeticID valueobject.CosmeticID) (bool, error) {
	monitoring.Repository(logger.ModelUserCosmetic, "exists", map[string]any{
		"user_id":     userID,
		"cosmetic_id": cosmeticID,
	})

	var exists bool
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM user_cosmetics
			WHERE user_id = $1 AND cosmetic_id = $2
		)`, userID.String(), uint64(cosmeticID),
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error verificando inventario: %w", err)
	}
	return exists, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID valueobject.UserID) ([]usercosmeticdomain.UserCosmetic, error) {
	monitoring.Repository(logger.ModelUserCosmetic, "find_by_user_id", map[string]any{"user_id": userID})

	rows, err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryContext(ctx, `
		SELECT user_id, cosmetic_id, purchased_at
		FROM user_cosmetics
		WHERE user_id = $1
		ORDER BY purchased_at DESC`, userID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("error listando inventario: %w", err)
	}
	defer rows.Close()

	var models []usercosmeticdata.Model
	for rows.Next() {
		var model usercosmeticdata.Model
		if err := rows.Scan(&model.UserID, &model.CosmeticID, &model.PurchasedAt); err != nil {
			return nil, fmt.Errorf("error leyendo inventario: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando inventario: %w", err)
	}
	return usercosmeticdata.ToDomainList(models), nil
}
