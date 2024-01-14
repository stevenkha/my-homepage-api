package routes

import (
	"os"
	"time"

	log "github.com/ccpaging/log4go"
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
		log.Debug("Error loading .env file: %s", err)
		return err
	}

	user := os.Getenv("user")
	if user == "" {
		log.Debug("Could not read user")
		return err
	}

	pass := os.Getenv("pass")
	if pass == "" {
		return err
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Debug("Could not initialize selenium")
		return err
	}
	defer wd.Quit()

	err = wd.Get(url)
	if err != nil {
		log.Debug("Could not get url")
		return err
	}

	login, err := wd.FindElement(selenium.ByClassName, "btn-login")
	if err != nil {
		log.Debug("could not find login button")
		return err
	}

	time.Sleep(2 * time.Second)

	userField, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		log.Debug("Could not find user field")
		return err
	}

	return nil
}
