package user

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "fragments/internal/core/models"
)

type Repository struct {
    pg *pgxpool.Pool
}

func New(pg *pgxpool.Pool) *Repository {
    return &Repository{pg: pg}
}

func (r Repository) Init() error {
    ctx := context.Background()
    query := ` CREATE TABLE users (
        id           text  unique not null,
        name         text  not null,
        email        text  not null,
        role         int   not null,
        passwordhash bytea not null
    )`
    _, err := r.pg.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("unable to create table: %w", err)
    }
    return nil
}

func (r Repository) Drop() error {
    ctx := context.Background()
    query := ` DROP TABLE users`
    _, err := r.pg.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("unable to drop table: %w", err)
    }
    return nil
}

func (r Repository) Reset() error {
    err := r.Drop()
    if err != nil {
        return err
    }
    return r.Init()
}

func (r Repository) Create(ctx context.Context, nu models.User) (models.User, error) {
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

func (r Repository) GetById(ctx context.Context, id string) (models.User, error) {
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

func (r Repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
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
