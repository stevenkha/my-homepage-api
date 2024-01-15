package routes

import (
	"fmt"
	"io/ioutil"
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
		log.Fatal("Could not read pass")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  "sessionid",
		Value: "1uyessrm638tbntlstlyadmpk2p9ohu9",
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
