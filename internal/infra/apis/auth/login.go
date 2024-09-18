package auth

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"

    "fragments/internal/core/models"
    "fragments/internal/infra/httpserver"
)

type LoginInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginOutput struct {
    AuthToken string `json:"authtoken"`
}

func makeToken(u models.User, secret []byte) (string, error) {
    now := time.Now()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":       u.Email,
        "iss":       "go-fit",
        "aud":       u.Role,
        "exp":       now.Add(time.Hour).Unix(),
        "iat":       now.Unix(),
        "user.Id":   u.Id,
        "user.Role": u.Role,
        "user.Name": u.Name,
    })
    tokenString, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
func (a AuthApi) handleLogin(c *gin.Context) {
    var input LoginInput
    err := c.ShouldBindJSON(&input)
    if err != nil {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }
    u, err := a.userService.AuthenticateByEmail(c, input.Email, input.Password)
    if err != nil {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }
    authToken, err := makeToken(u, []byte(a.secret))
    if err != nil {
        httpserver.HandleError(c, http.StatusBadRequest, "bad request", err)
        return
    }
    //c.SetSameSite(http.SameSiteLaxMode)
    //c.SetCookie("authtoken", authToken, 3600 * 24 * 30, "", "", false, true)
    c.JSON(http.StatusOK, LoginOutput{AuthToken: authToken})
}
