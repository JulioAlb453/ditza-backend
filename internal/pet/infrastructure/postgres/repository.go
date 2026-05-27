package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ditza/internal/pet/data"
	petdomain "ditza/internal/pet/domain"
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

func (r *Repository) Create(ctx context.Context, entity *petdomain.Pet) error {
	monitoring.Repository(logger.ModelPet, "create", map[string]any{"user_id": entity.UserID})

	model := data.ToModel(*entity)
	_, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		INSERT INTO pets (
			user_id, name, level, xp, mood,
			equipped_hat_id, equipped_shirt_id, equipped_background_id, equipped_accessory_id,
			last_interaction_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		model.UserID,
		model.Name,
		model.Level,
		model.XP,
		model.Mood,
		model.EquippedHatID,
		model.EquippedShirtID,
		model.EquippedBackgroundID,
		model.EquippedAccessoryID,
		model.LastInteractionAt,
		model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creando mascota: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, entity *petdomain.Pet) error {
	monitoring.Repository(logger.ModelPet, "update", map[string]any{"user_id": entity.UserID})

	model := data.ToModel(*entity)
	result, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE pets
		SET name = $2, level = $3, xp = $4, mood = $5,
		    equipped_hat_id = $6, equipped_shirt_id = $7,
		    equipped_background_id = $8, equipped_accessory_id = $9,
		    last_interaction_at = $10, updated_at = $11
		WHERE user_id = $1`,
		model.UserID,
		model.Name,
		model.Level,
		model.XP,
		model.Mood,
		model.EquippedHatID,
		model.EquippedShirtID,
		model.EquippedBackgroundID,
		model.EquippedAccessoryID,
		model.LastInteractionAt,
		model.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error actualizando mascota: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo verificar la actualización de mascota: %w", err)
	}
	if rowsAffected == 0 {
		return domainerror.New("PET_NOT_FOUND", "mascota no encontrada", domainerror.ErrNotFound)
	}
	return nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*petdomain.Pet, error) {
	monitoring.Repository(logger.ModelPet, "find_by_user_id", map[string]any{"user_id": userID})

	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT user_id, name, level, xp, mood,
		       equipped_hat_id, equipped_shirt_id, equipped_background_id, equipped_accessory_id,
		       last_interaction_at, updated_at
		FROM pets
		WHERE user_id = $1`, userID.String(),
	).Scan(
		&model.UserID,
		&model.Name,
		&model.Level,
		&model.XP,
		&model.Mood,
		&model.EquippedHatID,
		&model.EquippedShirtID,
		&model.EquippedBackgroundID,
		&model.EquippedAccessoryID,
		&model.LastInteractionAt,
		&model.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando mascota: %w", err)
	}

	entity := data.ToDomain(model)
	return &entity, nil
}
