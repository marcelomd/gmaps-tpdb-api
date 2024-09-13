package main

import (
	"log/slog"

	"fragments/internal/core/logger"
	api_auth "fragments/internal/infra/apis/auth"
	api_user "fragments/internal/infra/apis/user"
	"fragments/internal/infra/httpserver"
	"fragments/internal/infra/repositories/sqlite"
	repositories_user "fragments/internal/infra/repositories/sqlite/user"
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
	sqlite, err := sqlite.NewSqliteConnection("file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	userRepo := repositories_user.NewSqliteUserRepository(sqlite)
	userRepo.Init()

	// ----- Services
	userSvc := services_user.NewUserService(userRepo)

	// ----- Controllers
	authApi := api_auth.NewAuthApi(cfg.Secret, userSvc)
	userApi := api_user.NewUserApi(userSvc)

	// ----- Go
	server := httpserver.NewServer(cfg.Address, authApi.AuthMiddleware, authApi.RoleMiddleware)
	authApi.Register("/", server)
	userApi.Register("/user", server)

	server.Run()
}
