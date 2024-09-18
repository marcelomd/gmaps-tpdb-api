package main

import (
    "log/slog"

    "fragments/internal/core/logger"
    api_auth "fragments/internal/infra/apis/auth"
    api_user "fragments/internal/infra/apis/user"
    "fragments/internal/infra/httpserver"
    "fragments/internal/infra/repositories/postgres"
    repositories_user "fragments/internal/infra/repositories/postgres/user"
    services_user "fragments/internal/services/user"
)

func main() {
    logger.Init()

    cfg, err := loadConfig()
    if err != nil {
        panic(err)
    }
    slog.Info("Running config", slog.Any("config", cfg))

    // ----- Repositories
    pg, err := postgres.NewConnection(cfg.PGUrl)
    if err != nil {
        panic(err)
    }
    err = postgres.RunMigrations(cfg.PGUrl)
    if err != nil {
        panic(err)
    }

    userRepo := repositories_user.New(pg)
    _ = userRepo.Init()

    // ----- Services
    userSvc := services_user.New(userRepo)

    // ----- Controllers
    authApi := api_auth.New(cfg.Secret, userSvc)
    userApi := api_user.New(userSvc)

    // ----- Go!
    server := httpserver.New(cfg.Address, authApi.AuthMiddleware, authApi.RoleMiddleware)
    _ = authApi.Register("/", server)
    _ = userApi.Register("/user", server)

    _ = server.Run()
}
