package controllers_v1

import (
	"authService/config"
	services "authService/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

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

	go func() {
		account:= services.NewAccountService[any]()
		wg := sync.WaitGroup{}
		wg.Add(1)
		result, err := account.CreateOAuthAccount(user,&wg)
		if err != nil {
			log.Printf("Error creating OAuth account: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating account"})
			return
		}
		fmt.Println(result)
		c.Redirect(301,os.Getenv("CLIENT_SIDE"))
		wg.Wait()
	}()
}

