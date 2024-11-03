package config

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func LoadGoth() {
    // Load the session secret from the environment variable
    sessionSecret := os.Getenv("SESSION_SECRET")
    store := sessions.NewCookieStore([]byte(sessionSecret))
    gothic.Store = store

    // Configure Google provider
    googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
    googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    googleCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

    // Configure Microsoft provider
    microsoftClientId := os.Getenv("MICROSOFT_CLIENT_ID")
    microsoftClientSecret := os.Getenv("MICROSOFT_CLIENT_SECRET")
    microsoftCallbackURL := os.Getenv("MICROSOFT_CALLBACK_URL")

    // Configure Facebook provider
    facebookClientId := os.Getenv("FACEBOOK_CLIENT_ID")
    facebookClientSecret := os.Getenv("FACEBOOK_CLIENT_SECRET")
    facebookCallbackURL := os.Getenv("FACEBOOK_CALLBACK_URL")

    // Configure GitHub provider
    githubClientId := os.Getenv("GITHUB_CLIENT_ID")
    githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
    githubCallbackURL := os.Getenv("GITHUB_CALLBACK_URL")

    // Use configured providers with goth
    goth.UseProviders(
        google.New(
            googleClientId,
            googleClientSecret,
            googleCallbackURL,
            "email",
            "profile",
        ),
        microsoftonline.New(
            microsoftClientId,
            microsoftClientSecret,
            microsoftCallbackURL,
            "email",
            "profile",
            "offline_access",
            "openid",
            "User.Read",
            
        ),
        facebook.New(
            facebookClientId,
            facebookClientSecret,
            facebookCallbackURL,
            "email",
            "public_profile",
        ),
        github.New(
            githubClientId,
            githubClientSecret,
            githubCallbackURL,
            "user",
            "repo",
        ),
    )
}

