package postgres

import (
	"context"
	"database/sql"
)

type ctxKey struct{}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

func ExecutorFromContext(ctx context.Context, db *sql.DB) Executor {
	if tx, ok := ctx.Value(ctxKey{}).(*sql.Tx); ok && tx != nil {
		return tx
	}
	return db
}
