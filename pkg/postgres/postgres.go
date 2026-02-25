package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return pool, err
	}

	if err := pool.Ping(ctx); err != nil {
		return pool, err
	}

	return pool, nil
}
