package configs

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func LoadConfig() {
	// Oauth configuration for Google
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     "73653467917-vdn1egbp4pbchd5d1pgqtaqqljnuebhe.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/google_callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}
