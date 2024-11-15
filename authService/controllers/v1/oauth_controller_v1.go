package controllers_v1

import (
	"authService/config"
	"authService/internal/models"
	"authService/internal/repository"
	services "authService/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func init() {
	config.LoadGoth()
}

func BeginAuthHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackHandler handles the OAuth provider's callback after authentication
func CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication failed"})
		return
	}
	
	result := make(chan repository.ResultChan[models.User])
	go services.NewAccountService().CreateOAuthAccount(user, &result)

	select {
	case result := <-result:
		if result.Error != nil {
			log.Printf("User creation error: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		// Successfully created user, redirecting
		c.Redirect(http.StatusMovedPermanently, config.GetClientSide())
		return

	case <-time.After(30 * time.Second): // Timeout after 30 seconds
		log.Println("Timeout waiting for user creation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request timeout"})
		return
	}
}
