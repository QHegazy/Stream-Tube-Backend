package controllers_v1

import (
	"authService/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Define a custom claims structure
type Claims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}



// Auth generates an access token and a refresh token
func Auth(c *gin.Context) {
	accessToken, err := utils.GenerateToken("username", 15*time.Minute)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}
	c.JSON(200, gin.H{"access_token": accessToken})
}
