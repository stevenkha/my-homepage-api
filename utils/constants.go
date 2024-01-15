package utils

import (
	"os"

	log "github.com/ccpaging/log4go"
	"github.com/joho/godotenv"
)

const (
	AnimeUrl       = "https://yugenanime.tv/mylist/"
	AnimeListClass = "list-entries"

	ManagaUrl = "https://manganato.com/bookmark"
)

func GetEnvValues() (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Debug("Error loading .env file: %s", err)
		return "", "", err
	}

	cookieName := os.Getenv("animeCookieName")
	if cookieName == "" {
		log.Debug("Could not read cookie name")
		return "", "", err
	}

	cookieValue := os.Getenv("animeCookieValue")
	if cookieValue == "" {
		log.Fatal("Could not read cookie value")
	}

	return cookieName, cookieValue, nil
}
