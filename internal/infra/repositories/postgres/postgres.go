package postgres

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(source string) (*pgxpool.Pool, error) {
    ctx := context.Background()
    return pgxpool.New(ctx, source)
}
