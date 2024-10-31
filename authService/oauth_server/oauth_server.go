package oauthserver

import (
	v1 "authService/controllers/v1"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)


func OAuth() {
	r := mux.NewRouter()

	r.HandleFunc("/auth/{provider}", v1.BeginAuthHandler)
	r.HandleFunc("/auth/{provider}/callback", v1.CallbackHandler)

	port := os.Getenv("PORT")
	fmt.Printf("Starting server on https://localhost:%s...\n", port)

	err := http.ListenAndServeTLS(":"+port, "localhost+2.pem", "localhost+2-key.pem", r)
	if err != nil {
		log.Fatal(err)
	}
}

