package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"authService/config"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

// init sets up goth and loads OAuth provider configurations
func init() {
	config.LoadGoth() // Initializes Goth providers from environment variables
}

// main initializes the server and routes
func main() {
	r := mux.NewRouter()

	// OAuth routes
	r.HandleFunc("/auth/{provider}", beginAuthHandler)
	r.HandleFunc("/auth/{provider}/callback", callbackHandler)

	// Start server with HTTPS
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on https://localhost:%s...\n", port)

	// Use your self-signed certificate and key
	err := http.ListenAndServeTLS(":"+port, "localhost+2.pem", "localhost+2-key.pem", r)
	if err != nil {
		log.Fatal(err)
	}
}

// beginAuthHandler starts the OAuth flow by redirecting to the provider
func beginAuthHandler(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]
	q := r.URL.Query()
	q.Set("provider", provider)
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

// callbackHandler handles the OAuth provider's callback after authentication
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    // Explicitly set the provider to microsoftonline for Microsoft
    provider := mux.Vars(r)["provider"]
    q := r.URL.Query()
    q.Set("provider", provider)
    r.URL.RawQuery = q.Encode()

    user, err := gothic.CompleteUserAuth(w, r)
    if err != nil {
        log.Printf("Authentication error: %v", err)
        http.Error(w, "Authentication failed", http.StatusBadRequest)
        return
    }

    // Marshal user data to JSON
    userData, err := json.Marshal(user)
    if err != nil {
        http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
        return
    }

    // Set Content-Type and return JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(userData)
}
