package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetGoogleClientID() string {
	return os.Getenv("GOOGLE_CLIENT_ID")
}

func GetGoogleClientSecret() string {
	return os.Getenv("GOOGLE_CLIENT_SECRET")
}

func GetGoogleRedirectURI() string {
	return os.Getenv("GOOGLE_REDIRECT_URI")
}
