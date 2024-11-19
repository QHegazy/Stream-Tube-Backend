package oauth_server

import (
	v1 "authService/controllers/v1"
	"authService/middlewares"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func OAuth() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())
	// r.Use(middlewares.SecurityMiddleware(os.Getenv("HOST")))
	r.GET("/auth/:provider", v1.BeginAuthHandler)
	r.GET("/auth/:provider/callback", v1.CallbackHandler)
	r.POST("auth/register", v1.Auth)
	// r.GET("/h", func(c *gin.Context) {
	// 	c.SetCookie("refreshToken", "GG", 0, "/", ".hellogg.tv", true, false) // For subdomains of hellogg.tv
	// 	// or just		
	// 	c.Redirect(http.StatusFound, "https://hellogg.tv:4200")
	// })
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if not set
	}
	r.Run(":"+port)
	fmt.Printf("Starting server on http://localhost:%s...\n", port)
	
}
