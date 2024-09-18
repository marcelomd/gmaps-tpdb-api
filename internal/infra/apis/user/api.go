package user

import (
    "fragments/internal/core/models"
    "github.com/gin-gonic/gin"

    "fragments/internal/core/interfaces"
    "fragments/internal/infra/httpserver"
)

type UserApi struct {
    userService interfaces.UserService
}

func New(userService interfaces.UserService) *UserApi {
    return &UserApi{userService: userService}
}

func (a UserApi) Register(prefix string, server *httpserver.Server, middlewares ...gin.HandlerFunc) error {
    rg := server.Router.Group(prefix)
    for _, mw := range middlewares {
        rg.Use(mw)
    }
    rg.POST("/register", a.handleRegister)
    rg.GET("/me", server.RequireAuth(), a.handleGetMe)
    rg.POST("/", server.RequireAuth(), server.RequireRole(models.AdminRole), a.handleCreate)
    rg.GET("/:id", server.RequireAuth(), server.RequireRole(models.AdminRole), a.handleGetById)
    return nil
}
