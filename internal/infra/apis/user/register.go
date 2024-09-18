package user

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "fragments/internal/infra/httpserver"
    "fragments/internal/core/models"
)

type RegisterInput struct {
    Name      string `json:"name"`
    Email     string `json:"email"`
    Password1 string `json:"password1"`
    Password2 string `json:"password2"`
}

type RegisterOutput struct {
    Id    string `json:"id"`
    Name  string `json:"name"`
    Role  string `json:"role"`
    Email string `json:"email"`
}

func RegisterInputIsValid(input RegisterInput) bool {
    if input.Name == "" {
        return false
    }
    if input.Email == "" {
        return false
    }
    if input.Password1 == "" || input.Password2 == "" {
        return false
    }
    if input.Password1 != input.Password2 {
        return false
    }
    return true
}

func (a UserApi) handleRegister(c *gin.Context) {
    var input RegisterInput
    err := c.ShouldBindJSON(&input)
    if err != nil {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }

    if !RegisterInputIsValid(input) {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }

    newUser := models.NewUser{
        Name:     input.Name,
        Role:     models.UserRole,
        Email:    input.Email,
        Password: input.Password1,
    }
    result, err := a.userService.Create(c.Request.Context(), newUser)
    if err != nil {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }

    output := RegisterOutput{
        Id:    result.Id,
        Name:  result.Name,
        Email: result.Email,
        Role:  result.Role.String(),
    }
    c.JSON(http.StatusOK, output)
}
