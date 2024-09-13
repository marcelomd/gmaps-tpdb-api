package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"

	"fragments/internal/core/models"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	Id            string    `bun:"id,pk,type:uuid"`
	Created       time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Updated       time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Name          string
	Email         string
	Role          string
	PasswordHash  []byte `bun:"passwordhash"`
}

func toDb(u models.User) User {
	return User{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		Role:         u.Role,
		PasswordHash: u.PasswordHash,
	}
}

func fromDb(u User) models.User {
	return models.User{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		Role:         u.Role,
		PasswordHash: u.PasswordHash,
	}
}

type SqliteUserRepository struct {
	sqlite *sql.DB
	db     *bun.DB
}

func NewSqliteUserRepository(sqlite *sql.DB) *SqliteUserRepository {
	db := bun.NewDB(sqlite, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true))) // debug
	return &SqliteUserRepository{sqlite: sqlite, db: db}
}

func (r SqliteUserRepository) Init() error {
	c := context.Background()
	_, err := r.db.NewCreateTable().Model((*User)(nil)).Exec(c)
	if err != nil {
		return err
	}
	return nil
}
func (r SqliteUserRepository) Drop() error {
	c := context.Background()
	_, err := r.db.NewDropTable().Model((*User)(nil)).IfExists().Exec(c)
	if err != nil {
		return err
	}
	return nil
}
func (r SqliteUserRepository) Reset() error {
	err := r.Drop()
	if err != nil {
		return err
	}
	return r.Init()
}

func (r SqliteUserRepository) Create(ctx context.Context, nu models.User) (models.User, error) {
	user := toDb(nu)
	_, err := r.db.NewInsert().Model(&user).Returning("*").Exec(ctx)
	return fromDb(user), err
}

func (r SqliteUserRepository) GetById(ctx context.Context, id string) (models.User, error) {
	user := User{}
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	return fromDb(user), err
}

func (r SqliteUserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	user := User{}
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	return fromDb(user), err
}
