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

    userRepo := repositories_user.NewUserRepository(pg)
    _ = userRepo.Init()

    // ----- Services
    userSvc := services_user.NewUserService(userRepo)

    // ----- Controllers
    authApi := api_auth.NewAuthApi(cfg.Secret, userSvc)
    userApi := api_user.NewUserApi(userSvc)

    // ----- Go
    server := httpserver.NewServer(cfg.Address, authApi.AuthMiddleware, authApi.RoleMiddleware)
    _ = authApi.Register("/", server)
    _ = userApi.Register("/user", server)

    _ = server.Run()
}
