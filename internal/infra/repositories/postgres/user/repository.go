package user

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "fragments/internal/core/models"
)

type UserRepository struct {
    pg *pgxpool.Pool
}

func NewUserRepository(pg *pgxpool.Pool) *UserRepository {
    return &UserRepository{pg: pg}
}

func (r UserRepository) Init() error {
    ctx := context.Background()
    query := ` CREATE TABLE users (
        id           text  unique not null,
        name         text  not null,
        email        text  not null,
        role         text  not null,
        passwordhash bytea not null
    )`
    _, err := r.pg.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("unable to create table: %w", err)
    }
    return nil
}

func (r UserRepository) Drop() error {
    ctx := context.Background()
    query := ` DROP TABLE users`
    _, err := r.pg.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("unable to drop table: %w", err)
    }
    return nil
}

func (r UserRepository) Reset() error {
    err := r.Drop()
    if err != nil {
        return err
    }
    return r.Init()
}

func (r UserRepository) Create(ctx context.Context, nu models.User) (models.User, error) {
    query := `INSERT INTO users (id, name, email, role, passwordhash) VALUES (@id, @name, @email, @role, @passwordhash)`
    args := pgx.NamedArgs{
        "id":           nu.Id,
        "name":         nu.Name,
        "email":        nu.Email,
        "role":         nu.Role,
        "passwordhash": nu.PasswordHash,
    }
    _, err := r.pg.Exec(ctx, query, args)
    if err != nil {
        return nu, fmt.Errorf("unable to insert row: %w", err)
    }
    return nu, nil
}

func (r UserRepository) GetById(ctx context.Context, id string) (models.User, error) {
    query := `SELECT id, name, email, role, passwordhash FROM users WHERE id = @id`
    args := pgx.NamedArgs{"id": id}
    row := r.pg.QueryRow(ctx, query, args)
    user := models.User{}
    err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.PasswordHash)
    switch err {
    case nil:
        return user, err
    case pgx.ErrNoRows:
        return user, fmt.Errorf("not found")
    default:
        return user, fmt.Errorf("unable to get user %s: %w", id, err)
    }
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
    query := `SELECT id, name, email, role, passwordhash FROM users WHERE email = @email`
    args := pgx.NamedArgs{"email": email}
    row := r.pg.QueryRow(ctx, query, args)
    user := models.User{}
    err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.PasswordHash)
    switch err {
    case nil:
        return user, err
    case pgx.ErrNoRows:
        return user, fmt.Errorf("not found")
    default:
        return user, fmt.Errorf("unable to get user %s: %w", email, err)
    }
}
