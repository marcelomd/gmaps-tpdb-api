package user

import (
    "log/slog"
    "net/http"

    "github.com/gin-gonic/gin"

    "fragments/internal/core/models"
    "fragments/internal/infra/httpserver"
)

func (a UserApi) handleGetMe(c *gin.Context) {
    authFromCtx, exists := c.Get("userAuth")
    if !exists {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", nil)
        return
    }

    auth := authFromCtx.(models.AuthData)
    slog.Info("auth data", slog.Any("auth", auth))

    result, err := a.userService.GetById(c.Request.Context(), auth.Id)
    if err != nil {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }

    output := GetByIdOutput{
        Id:    result.Id,
        Name:  result.Name,
        Email: result.Email,
        Role:  result.Role,
    }
    c.JSON(http.StatusOK, output)
}
