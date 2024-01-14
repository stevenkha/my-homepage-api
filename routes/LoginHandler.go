package routes

import (
	"fmt"
	"log"
	"os"
	"time"

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

	err = wd.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	login, err := wd.FindElement(selenium.ByClassName, "btn-login")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	return nil
}
