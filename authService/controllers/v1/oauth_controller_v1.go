package v1

import (
	"authService/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func init() {
	config.LoadGoth() 
}


func BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["provider"]
	q := r.URL.Query()
	q.Set("provider", provider)
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

// callbackHandler handles the OAuth provider's callback after authentication
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
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
    fmt.Println(json.Marshal(user))
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
    
}
