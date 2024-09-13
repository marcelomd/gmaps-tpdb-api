package user

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"fragments/internal/core/models"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	auth := authFromCtx.(models.AuthData)
	slog.Info("auth data", slog.Any("auth", auth))

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	result, err := a.userService.GetById(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
