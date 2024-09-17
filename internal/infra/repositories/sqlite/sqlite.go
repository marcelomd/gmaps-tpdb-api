package sqlite

import (
    "database/sql"

    "github.com/uptrace/bun/driver/sqliteshim"
)

func NewConnection(source string) (*sql.DB, error) {
    return sql.Open(sqliteshim.ShimName, source)
}
