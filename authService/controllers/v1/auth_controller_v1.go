package controllers_v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key used for signing the JWT
var secretKey = []byte("your_secret_key")

// Define a custom claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func auth(w http.ResponseWriter, r *http.Request) {
	// In a real-world application, you'd verify the user credentials first
	// For this example, we'll assume the user is valid if they provide a username
	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Create the JWT claims, which includes the username and expiration time
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Set expiration time to 24 hours
			Issuer:    "my-app", // Issuer of the token
		},
	}

	// Create the token using the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Error signing the token", http.StatusInternalServerError)
		return
	}

	// Send the JWT back in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"token": "%s"}`, signedToken)
}
