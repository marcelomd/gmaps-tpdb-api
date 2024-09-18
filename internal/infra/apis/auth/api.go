package auth

import (
    "github.com/gin-gonic/gin"

    "fragments/internal/core/interfaces"
    "fragments/internal/infra/httpserver"
)

type AuthApi struct {
    secret      string
    userService interfaces.UserService
}

func New(secret string, userService interfaces.UserService) *AuthApi {
    return &AuthApi{secret: secret, userService: userService}
}

func (a AuthApi) Register(prefix string, server *httpserver.Server, middlewares ...gin.HandlerFunc) error {
    rg := server.Router.Group(prefix)
    for _, mw := range middlewares {
        rg.Use(mw)
    }
    rg.POST("/login", a.handleLogin)
    rg.POST("/refresh", server.RequireAuth(), a.handleRefresh)
    return nil
}

func (a AuthApi) handleRefresh(c *gin.Context) {

}
