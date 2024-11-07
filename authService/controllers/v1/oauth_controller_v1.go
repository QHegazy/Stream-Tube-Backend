package controllers_v1

import (
	"authService/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	fmt.Println(json.Marshal(user))
	c.JSON(http.StatusOK, user)
}
