package configs

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func LoadConfig() {
	// Oauth configuration for Google
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/google_callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}
