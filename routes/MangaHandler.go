package routes

import (
	"my-homepage-api/utils"
	"net/http"
	"net/http/cookiejar"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
)

func MangaHandler(c *gin.Context) error {
	// TODO: Initializing cookies and client is redundant for both Manga and Anime
	// Clean it up later

	log.Info("Loading env...")

	cookieName, cookieValue, err := utils.GetEnvValues("mangaCookieName", "mangaCookieValue")
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

	req, err := http.NewRequest("GET", utils.ManagaUrl, nil)
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
