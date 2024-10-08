package user

import (
    "log/slog"
    "net/http"

    "github.com/gin-gonic/gin"

    "fragments/internal/core/models"
    "fragments/internal/infra/httpserver"
)

type GetByIdOutput struct {
    Id    string `json:"id"`
    Name  string `json:"name"`
    Role  string `json:"role"`
    Email string `json:"email"`
}

func (a UserApi) handleGetById(c *gin.Context) {
    authFromCtx, exists := c.Get("userAuth")
    if !exists {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", nil)
        return
    }

    auth := authFromCtx.(models.AuthData)
    slog.Info("auth data", slog.Any("auth", auth))

    id := c.Param("id")
    if id == "" {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", nil)
        return
    }

    result, err := a.userService.GetById(c.Request.Context(), id)
    if err != nil {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }

    output := GetByIdOutput{
        Id:    result.Id,
        Name:  result.Name,
        Email: result.Email,
        Role:  result.Role.String(),
    }
    c.JSON(http.StatusOK, output)
}
