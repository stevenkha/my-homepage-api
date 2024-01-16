package utils

import (
	"os"

	log "github.com/ccpaging/log4go"
	"github.com/joho/godotenv"
)

const (
	AnimeUrl       = "https://yugenanime.tv/mylist/"
	AnimeListClass = "list-entries"

	MangaUrl = "https://user.mngusr.com/bookmark_get_list_full"
)

func GetEnvValues(name string, value string) (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Debug("Error loading .env file: %s", err)
		return "", "", err
	}

	cookieName := os.Getenv(name)
	if cookieName == "" {
		log.Debug("Could not read cookie name")
		return "", "", err
	}

	cookieValue := os.Getenv(value)
	if cookieValue == "" {
		log.Fatal("Could not read cookie value")
	}

	return cookieName, cookieValue, nil
}
