package user

import (
	"github.com/gin-gonic/gin"

	"fragments/internal/core/interfaces"
	"fragments/internal/infra/httpserver"
)

type UserApi struct {
	userService interfaces.UserService
}

func NewUserApi(userService interfaces.UserService) *UserApi {
	return &UserApi{userService: userService}
}

func (a UserApi) Register(prefix string, server *httpserver.Server, middlewares ...gin.HandlerFunc) error {
	rg := server.Router.Group(prefix)
	for _, mw := range middlewares {
		rg.Use(mw)
	}
	rg.POST("/", a.handleCreate)
	rg.GET("/me", server.RequireAuth(), a.handleGetMe)
	rg.GET("/:id", server.RequireAuth(), server.RequireRole("staff"), a.handleGetById)
	return nil
}
