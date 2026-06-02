package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ditza/internal/friendship/data"
	friendshipdomain "ditza/internal/friendship/domain"
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

func (r *Repository) Create(ctx context.Context, entity *friendshipdomain.Friendship) error {
	monitoring.Repository(logger.ModelFriendship, "create", map[string]any{
		"requester_id": entity.RequesterID,
		"addressee_id": entity.AddresseeID,
	})

	model := data.ToModel(*entity)
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		INSERT INTO friendships (requester_id, addressee_id, status, created_at, responded_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		model.RequesterID,
		model.AddresseeID,
		model.Status,
		model.CreatedAt,
		model.RespondedAt,
	).Scan(&model.ID)
	if err != nil {
		return fmt.Errorf("error creando amistad: %w", err)
	}

	entity.ID = valueobject.FriendshipID(model.ID)
	return nil
}

func (r *Repository) Update(ctx context.Context, entity *friendshipdomain.Friendship) error {
	monitoring.Repository(logger.ModelFriendship, "update", map[string]any{"friendship_id": entity.ID})

	model := data.ToModel(*entity)
	result, err := sharedpostgres.ExecutorFromContext(ctx, r.db).ExecContext(ctx, `
		UPDATE friendships
		SET status = $2, responded_at = $3
		WHERE id = $1`,
		model.ID,
		model.Status,
		model.RespondedAt,
	)
	if err != nil {
		return fmt.Errorf("error actualizando amistad: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("no se pudo verificar la actualización de amistad: %w", err)
	}
	if rowsAffected == 0 {
		return domainerror.New("FRIENDSHIP_NOT_FOUND", "solicitud de amistad no encontrada", domainerror.ErrNotFound)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id valueobject.FriendshipID) (*friendshipdomain.Friendship, error) {
	monitoring.Repository(logger.ModelFriendship, "find_by_id", map[string]any{"friendship_id": id})

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

func (r *Repository) FindBetweenUsers(ctx context.Context, userA, userB valueobject.UserID) (*friendshipdomain.Friendship, error) {
	monitoring.Repository(logger.ModelFriendship, "find_between_users", map[string]any{
		"user_a": userA,
		"user_b": userB,
	})

	model, err := r.queryBetweenUsers(ctx, userA, userB)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, nil
	}

	entity := data.ToDomain(*model)
	return &entity, nil
}

func (r *Repository) FindAcceptedByUserID(ctx context.Context, userID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	monitoring.Repository(logger.ModelFriendship, "find_accepted_by_user_id", map[string]any{"user_id": userID})

	return r.queryList(ctx, `
		WHERE status = 'accepted'
		  AND (requester_id = $1 OR addressee_id = $1)
		ORDER BY created_at DESC`, userID.String())
}

func (r *Repository) FindPendingByAddresseeID(ctx context.Context, addresseeID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	monitoring.Repository(logger.ModelFriendship, "find_pending_by_addressee_id", map[string]any{"addressee_id": addresseeID})

	return r.queryList(ctx, `
		WHERE status = 'pending' AND addressee_id = $1
		ORDER BY created_at DESC`, addresseeID.String())
}

func (r *Repository) queryOne(ctx context.Context, whereClause string, args ...any) (*data.Model, error) {
	query := fmt.Sprintf(`
		SELECT id, requester_id, addressee_id, status, created_at, responded_at
		FROM friendships %s`, whereClause)

	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, query, args...).Scan(
		&model.ID,
		&model.RequesterID,
		&model.AddresseeID,
		&model.Status,
		&model.CreatedAt,
		&model.RespondedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando amistad: %w", err)
	}
	return &model, nil
}

func (r *Repository) queryList(ctx context.Context, whereClause string, args ...any) ([]friendshipdomain.Friendship, error) {
	query := fmt.Sprintf(`
		SELECT id, requester_id, addressee_id, status, created_at, responded_at
		FROM friendships %s`, whereClause)

	rows, err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listando amistades: %w", err)
	}
	defer rows.Close()

	var models []data.Model
	for rows.Next() {
		var model data.Model
		if err := rows.Scan(
			&model.ID,
			&model.RequesterID,
			&model.AddresseeID,
			&model.Status,
			&model.CreatedAt,
			&model.RespondedAt,
		); err != nil {
			return nil, fmt.Errorf("error leyendo amistad: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando amistades: %w", err)
	}
	return data.ToDomainList(models), nil
}

func (r *Repository) queryBetweenUsers(ctx context.Context, userA, userB valueobject.UserID) (*data.Model, error) {
	var model data.Model
	err := sharedpostgres.ExecutorFromContext(ctx, r.db).QueryRowContext(ctx, `
		SELECT id, requester_id, addressee_id, status, created_at, responded_at
		FROM friendships
		WHERE LEAST(requester_id, addressee_id) = LEAST($1::uuid, $2::uuid)
		  AND GREATEST(requester_id, addressee_id) = GREATEST($1::uuid, $2::uuid)`,
		userA.String(), userB.String(),
	).Scan(
		&model.ID,
		&model.RequesterID,
		&model.AddresseeID,
		&model.Status,
		&model.CreatedAt,
		&model.RespondedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando amistad entre usuarios: %w", err)
	}
	return &model, nil
}
