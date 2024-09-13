package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"fragments/internal/core/interfaces"
	"fragments/internal/core/models"
	"fragments/internal/infra/httpserver"
)

type AuthApi struct {
	secret      string
	userService interfaces.UserService
}

func NewAuthApi(secret string, userService interfaces.UserService) *AuthApi {
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := a.userService.AuthenticateByEmail(c, input.Email, input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authToken, err := makeToken(u, []byte(a.secret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("authtoken", authToken, 3600 * 24 * 30, "", "", false, true)
	c.JSON(http.StatusOK, LoginOutput{AuthToken: authToken})
}

func (a AuthApi) handleRefresh(c *gin.Context) {

}
