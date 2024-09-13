package sqlite

import (
	"database/sql"

	"github.com/uptrace/bun/driver/sqliteshim"
)

func NewSqliteConnection(source string) (*sql.DB, error) {
	return sql.Open(sqliteshim.ShimName, source)
}
