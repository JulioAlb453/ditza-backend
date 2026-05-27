package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"ditza/internal/point-transaction/data"
	pointtransactiondomain "ditza/internal/point-transaction/domain"
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

func (r *Repository) Create(ctx context.Context, entity *pointtransactiondomain.PointTransaction) error {
	monitoring.Repository(logger.ModelPointTransaction, "create", map[string]any{
		"user_id": entity.UserID,
		"type":    entity.Type,
	})

	model := data.ToModel(*entity)
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		INSERT INTO point_transactions (
			user_id, type, coins_delta, season_delta, reference_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		model.UserID,
		model.Type,
		model.CoinsDelta,
		model.SeasonDelta,
		model.ReferenceID,
		model.CreatedAt,
	).Scan(&model.ID)
	if err != nil {
		return fmt.Errorf("error creando transacción de puntos: %w", err)
	}

	entity.ID = valueobject.PointTransactionID(model.ID)
	return nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID valueobject.UserID, limit int) ([]pointtransactiondomain.PointTransaction, error) {
	monitoring.Repository(logger.ModelPointTransaction, "find_by_user_id", map[string]any{
		"user_id": userID,
		"limit":   limit,
	})

	if limit <= 0 {
		limit = 50
	}

	rows, err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryContext(ctx, `
		SELECT id, user_id, type, coins_delta, season_delta, reference_id, created_at
		FROM point_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, userID.String(), limit,
	)
	if err != nil {
		return nil, fmt.Errorf("error listando transacciones: %w", err)
	}
	defer rows.Close()

	var models []data.Model
	for rows.Next() {
		var model data.Model
		if err := rows.Scan(
			&model.ID,
			&model.UserID,
			&model.Type,
			&model.CoinsDelta,
			&model.SeasonDelta,
			&model.ReferenceID,
			&model.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error leyendo transacción: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando transacciones: %w", err)
	}
	return data.ToDomainList(models), nil
}
