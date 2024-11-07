package oauthserver

import (
	v1 "authService/controllers/v1"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func OAuth() {
	r := gin.Default()

	// Routes for OAuth authentication
	r.GET("/auth/:provider", v1.BeginAuthHandler)
	r.GET("/auth/:provider/callback", v1.CallbackHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if not set
	}

	fmt.Printf("Starting server on https://localhost:%s...\n", port)

	err := r.RunTLS(":"+port, "localhost+2.pem", "localhost+2-key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
