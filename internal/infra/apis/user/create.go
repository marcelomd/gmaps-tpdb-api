package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fragments/internal/core/models"
)

type CreateInput struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateOutput struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

func CreateInputIsValid(input CreateInput) bool {
	return true
}

func (a UserApi) handleCreate(c *gin.Context) {
	var input CreateInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !CreateInputIsValid(input) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	newUser := models.NewUser{
		Name:     input.Name,
		Role:     input.Role,
		Email:    input.Email,
		Password: input.Password,
	}
	result, err := a.userService.Create(c.Request.Context(), newUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output := CreateOutput{
		Id:    result.Id,
		Name:  result.Name,
		Email: result.Email,
		Role:  result.Role,
	}
	c.JSON(http.StatusOK, output)
}
