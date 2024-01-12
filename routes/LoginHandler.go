package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	url = "https://aniwatch.to/user/notification?type=1"
)

func LoginHandler(c *gin.Context) error {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
		return err
	}

	token := os.Getenv("token")
	if token == "" {
		log.Fatal("Token is not set")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Could not process request")
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not get response")
		return err
	}
	defer resp.Body.Close()

	return nil
}
