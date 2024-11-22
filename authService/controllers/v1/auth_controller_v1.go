package controllers_v1

import (
	dto "authService/Dto"
	"net/http"

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
	var registerRequest dto.RegisterRequest
	err := c.ShouldBindJSON(&registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	 c.JSONP(200,registerRequest)
	 



}
