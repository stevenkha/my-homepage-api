package routes

import (
	"net/http"
	"net/http/cookiejar"
	"os"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	url = "https://yugenanime.tv/mylist/"
)

func LoginHandler(c *gin.Context) error {

	log.Info("Loading env...")

	cookieName, cookieValue, err := getEnvValues()
	if err != nil {
		log.Error(err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error("Could not initialize cookiejar")
		return err
	}

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Could not create request")
		return err
	}

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: cookieValue,
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getEnvValues() (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Debug("Error loading .env file: %s", err)
		return "", "", err
	}

	cookieName := os.Getenv("cookieName")
	if cookieName == "" {
		log.Debug("Could not read cookie name")
		return "", "", err
	}

	cookieValue := os.Getenv("cookieValue")
	if cookieValue == "" {
		log.Fatal("Could not read cookie value")
	}

	return cookieName, cookieValue, nil
}
