package oauth_server

import (
	v1 "authService/controllers/v1"
	"authService/middlewares"
	"authService/utils"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func OAuth() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.SecurityMiddleware(os.Getenv("HOST")))
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/:provider", v1.BeginAuthHandler)
		authGroup.GET("/:provider/callback", v1.CallbackHandler)
		authGroup.POST("/register", v1.Auth)
		authGroup.GET("/h", func(c *gin.Context) {
			codes := utils.GenerateRecoveryCodes(20)
			obj := map[string]interface{}{
				"code": codes,
			}
			c.JSONP(200, obj)
		})
	}
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}
	r.Run(":"+port)
	fmt.Printf("Starting server on http://localhost:%s...\n", port)
	
}
