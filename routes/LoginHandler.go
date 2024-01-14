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
	homeUrl = "https://aniwatch.to/home"
)

func LoginHandler(c *gin.Context) error {

	log.Info("Loading env...")

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

	log.Info("Initializng selenium...")

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Debug("Could not initialize selenium")
		return err
	}
	defer wd.Quit()

	log.Info("Getting url...")

	err = wd.Get(homeUrl)
	if err != nil {
		log.Debug("Could not get url")
		return err
	}

	log.Info("Finding login button...")

	login, err := wd.FindElement(selenium.ByCSSSelector, "[data-target='#modallogin']")
	if err != nil {
		log.Debug("Could not find login button")
		return err
	}

	if login.Click(); err != nil {
		log.Debug("Could not click on login button")
		return err
	}

	time.Sleep(5 * time.Second)

	log.Info("Finding user field...")

	userField, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		log.Debug("Could not find user field")
		return err
	}

	log.Info("Finding pass field...")

	passField, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		log.Debug("Could not find pass field")
		return err
	}

	log.Info("Sending keys...")

	if sendKeys(userField, passField, user, pass); err != nil {
		log.Debug("Could not fill user and pass fields")
		return err
	}

	log.Info("Finding submit button...")

	submitLogin, err := wd.FindElement(selenium.ByID, "btn-login")
	if err != nil {
		log.Debug("Could not find submit button for login")
		return err
	}

	log.Info("Clicking submit button...")

	if submitLogin.Click(); err != nil {
		log.Debug("Could not process submit login request")
		return err
	}

	time.Sleep(5 * time.Second)

	log.Info("LOGGED IN SUCCESSFULLY")

	return nil
}

func sendKeys(userField selenium.WebElement, passField selenium.WebElement, user string, pass string) error {
	err := userField.SendKeys(user)
	if err != nil {
		return err
	}

	err = passField.SendKeys(pass)
	if err != nil {
		return err
	}

	return nil
}
