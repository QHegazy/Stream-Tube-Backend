package controllers_v1

import (
	dto "authService/Dto"
	services "authService/internal/service"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}



func Auth(c *gin.Context) {
	var register dto.RegisterLocalUser
	err := c.ShouldBindJSON(&register)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go func() {
		account:= services.NewAccountService[string]()
		wg := sync.WaitGroup{}
		wg.Add(1)
		result, err := account.CreateLocalAuthAccount(&register,&wg)
		if err != nil {
			log.Printf("Error creating local user account: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating account"})
			return
		}
		wg.Wait()
		c.String(200,result)
	}()

}
