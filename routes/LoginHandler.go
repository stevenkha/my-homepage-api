package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
)

const (
	url = "https://aniwatch.to/home"
)

func LoginHandler(c *gin.Context) error {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
		return err
	}

	user := os.Getenv("user")
	if user == "" {
		log.Fatal("Could not read user")
	}

	pass := os.Getenv("pass")
	if pass == "" {
		log.Fatal("Could not read pass")
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal(err)
	}
	defer wd.Quit()

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
