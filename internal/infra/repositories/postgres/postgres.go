package postgres

import (
    "context"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewConnection(source string) (*pgxpool.Pool, error) {
    return pgxpool.New(context.Background(), source)
}
func RunMigrations(source string) error {
    s := strings.Replace(source, "postgres://", "pgx5://", 1)
    m, err := migrate.New("file://migrations", s)
    if err != nil {
        return err
    }
    err = m.Up()
    if err != nil {
        return err
    }
    return nil
}
