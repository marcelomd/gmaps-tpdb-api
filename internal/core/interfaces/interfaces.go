package interfaces

import (
    "context"

    "fragments/internal/core/models"
)

// ----- Services
type UserService interface {
    Create(ctx context.Context, nu models.NewUser) (models.User, error)
    GetById(ctx context.Context, id string) (models.User, error)
    AuthenticateByEmail(ctx context.Context, email string, password string) (models.User, error)
}

type CompoundService interface {
}

// ----- Repositories
type UserRepository interface {
    Create(ctx context.Context, nu models.User) (models.User, error)
    GetById(ctx context.Context, id string) (models.User, error)
    GetByEmail(ctx context.Context, email string) (models.User, error)
}

type CompoundRepository interface {
}
