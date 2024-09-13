package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"fragments/internal/core/models"
)

func tokenFromHeader(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if len(authorization) > 7 && strings.ToLower(authorization[0:6]) == "bearer" {
		return authorization[7:], nil
	}
	return "", fmt.Errorf("invalid authorization header")
}

func (a AuthApi) AuthMiddleware(c *gin.Context) {
	tokenString, err := tokenFromHeader(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	}
	token, _ := jwt.Parse(tokenString, keyFunc)
	if !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userAuth := models.AuthData{
		Id:   claims["user.Id"].(string),
		Role: claims["user.Role"].(string),
		Name: claims["user.Name"].(string),
	}
	c.Set("userAuth", userAuth)
	c.Next()
}

func (a AuthApi) RoleMiddleware(c *gin.Context, role string) {
	userAuth, exists := c.Get("userAuth")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if userAuth.(models.AuthData).Role != role {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}
