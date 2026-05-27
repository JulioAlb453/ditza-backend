package postgres

import (
	"context"
	"database/sql"
	"fmt"

	leaderboarddomain "ditza/internal/leaderboard/domain"
	valueobject "ditza/internal/shared/domain/value-object"
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetFriendRanking(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) ([]leaderboarddomain.FriendEntry, error) {
	monitoring.Repository(logger.ModelLeaderboard, "get_friend_ranking", map[string]any{
		"user_id":   userID,
		"season_id": seasonID,
	})

	rows, err := r.db.QueryContext(ctx, `
		WITH friend_ids AS (
			SELECT CASE
				WHEN requester_id = $1::uuid THEN addressee_id
				ELSE requester_id
			END AS participant_id
			FROM friendships
			WHERE status = 'accepted'
			  AND (requester_id = $1::uuid OR addressee_id = $1::uuid)
		),
		participants AS (
			SELECT $1::uuid AS participant_id
			UNION
			SELECT participant_id FROM friend_ids
		)
		SELECT u.id, u.alias, COALESCE(ss.points, 0) AS season_points
		FROM participants p
		JOIN users u ON u.id = p.participant_id
		LEFT JOIN season_scores ss ON ss.user_id = u.id AND ss.season_id = $2
		ORDER BY season_points DESC, u.alias ASC`,
		userID.String(), uint64(seasonID),
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ranking de amigos: %w", err)
	}
	defer rows.Close()

	var entries []leaderboarddomain.FriendEntry
	rank := 0
	var lastPoints *int

	for rows.Next() {
		var id string
		var alias string
		var points int
		if err := rows.Scan(&id, &alias, &points); err != nil {
			return nil, fmt.Errorf("error leyendo ranking: %w", err)
		}

		rank++
		if lastPoints != nil && points < *lastPoints {
			rank = len(entries) + 1
		}
		pointsCopy := points
		lastPoints = &pointsCopy

		parsedUserID := valueobject.UserID(id)
		entries = append(entries, leaderboarddomain.FriendEntry{
			UserID:        parsedUserID,
			Alias:         alias,
			SeasonPoints:  points,
			Rank:          rank,
			IsCurrentUser: parsedUserID == userID,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando ranking: %w", err)
	}
	if entries == nil {
		return []leaderboarddomain.FriendEntry{}, nil
	}
	return entries, nil
}
