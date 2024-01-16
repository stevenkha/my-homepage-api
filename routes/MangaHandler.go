package routes

import (
	"io/ioutil"
	"my-homepage-api/utils"
	"net/http"
	"net/url"
	"strings"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
)

func MangaHandler(c *gin.Context) error {

	cookieName, cookieValue, err := utils.GetEnvValues("bookmarkDataName", "bookmarkDataValue")
	if err != nil {
		return err
	}

	client := &http.Client{}

	formData := url.Values{}
	formData.Set(cookieName, cookieValue)
	formData.Set("bm_source", "manganato")
	formData.Set("out_type", "html")

	req, err := http.NewRequest("POST", utils.MangaUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Debug("Could not make post request to manganato api")
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug("Could not read post response body")
		return err
	}

	log.Debug(string(body))

	return nil
}
