package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	sharedpostgres "ditza/internal/shared/infrastructure/postgres"
	"ditza/internal/user/data"
	userdomain "ditza/internal/user/domain"

	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, entity *userdomain.User) error {
	monitoring.Repository(logger.ModelUser, "create", map[string]any{
		"user_id": entity.ID,
		"email":   entity.Email,
	})

	model := data.ToModel(*entity)
	_, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		INSERT INTO users (id, alias, email, password, coins, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		model.ID,
		model.Alias,
		model.Email,
		model.Password,
		model.Coins,
		model.CreatedAt,
		model.UpdatedAt,
	)
	if err != nil {
		return mapUserError("create", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, entity *userdomain.User) error {
	monitoring.Repository(logger.ModelUser, "update", map[string]any{"user_id": entity.ID})

	model := data.ToModel(*entity)
	result, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE users
		SET alias = $2, email = $3, password = $4, coins = $5, updated_at = $6
		WHERE id = $1`,
		model.ID,
		model.Alias,
		model.Email,
		model.Password,
		model.Coins,
		model.UpdatedAt,
	)
	if err != nil {
		return mapUserError("update", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo verificar la actualización del usuario: %w", err)
	}
	if rowsAffected == 0 {
		return domainerror.New("USER_NOT_FOUND", "usuario no encontrado", domainerror.ErrNotFound)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id valueobject.UserID) (*userdomain.User, error) {
	monitoring.Repository(logger.ModelUser, "find_by_id", map[string]any{"user_id": id})

	model, err := r.queryOne(ctx, "WHERE id = $1", id.String())
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}

	userEntity := data.ToDomain(*model)
	return &userEntity, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	monitoring.Repository(logger.ModelUser, "find_by_email", map[string]any{"email": email})

	model, err := r.queryOne(ctx, "WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}

	userEntity := data.ToDomain(*model)
	return &userEntity, nil
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	monitoring.Repository(logger.ModelUser, "exists_by_email", map[string]any{"email": email})

	var exists bool
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("no se pudo verificar si el correo existe: %w", err)
	}
	return exists, nil
}

func (r *Repository) queryOne(ctx context.Context, whereClause string, arg any) (*data.Model, error) {
	query := fmt.Sprintf(`
		SELECT id, alias, email, password, coins, created_at, updated_at
		FROM users
		%s`, whereClause)

	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, query, arg).Scan(
		&model.ID,
		&model.Alias,
		&model.Email,
		&model.Password,
		&model.Coins,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo consultar el usuario: %w", err)
	}
	return &model, nil
}

func mapUserError(operation string, err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return domainerror.New("USER_ALREADY_EXISTS", "el correo ya está registrado", domainerror.ErrInvalidInput)
	}
	return fmt.Errorf("error en repositorio de usuario (%s): %w", operation, err)
}
