package config

import (
	"os"

	"golang.org/x/oauth2"
)

func getOauthDev() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:65010/oauth/stripe_callback",
		ClientID:     os.Getenv("STRIPE_CLIENT_ID"),
		ClientSecret: os.Getenv("STRIPE_CLIENT_SECRET"),
		Scopes:       []string{"read_write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://connect.stripe.com/oauth/authorize",
			TokenURL: "https://connect.stripe.com/oauth/token",
		},
	}
}

func getProdcutionConfig() *oauth2.Config {

	return &oauth2.Config{
		RedirectURL:  os.Getenv("CALLBACK_URL"),
		ClientID:     os.Getenv("STRIPE_CLIENT_ID"),
		ClientSecret: os.Getenv("STRIPE_CLIENT_SECRET"),
		Scopes:       []string{"read_write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://connect.stripe.com/oauth/authorize",
			TokenURL: "https://connect.stripe.com/oauth/token",
		},
	}
}

// GetStripeOauthConfig returns the correct config based on the ENVIRONMENT
func GetStripeOauthConfig() *oauth2.Config {
	switch os.Getenv("ENVIRONMENT") {
	case "PRODUCTION":
	case "PRE-PRODUCTION":
		return getProdcutionConfig()
	default:
		return getOauthDev()
	}
	return nil
}
