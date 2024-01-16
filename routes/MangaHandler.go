package routes

import (
	"encoding/json"
	"io/ioutil"
	"my-homepage-api/utils"
	"net/http"
	"net/url"
	"strings"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type BookmarkResponse struct {
	Result string `json:"result"`
	Data   string `json:"data"`
}

type MangaInfo struct {
	Cover   string `json:"cover"`
	Title   string `json:"title"`
	Viewed  string `json:"viewed"`
	Current string `json:"current"`
}

func MangaHandler(c *gin.Context) error {

	dataName, dataValue, err := utils.GetEnvValues("bookmarkDataName", "bookmarkDataValue")
	if err != nil {
		return err
	}

	client := &http.Client{}

	formData := url.Values{}
	formData.Set(dataName, dataValue)
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

	var postPayload BookmarkResponse
	err = json.Unmarshal(body, &postPayload)
	if err != nil {
		return err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Debug("Error parsing html: ")
		return err
	}

	bookmarkListDiv := utils.GetListDiv(doc, "user-bookmark-content")
	if bookmarkListDiv == nil {
		log.Debug("Could not find bookmark list node")
	}

	mangas := utils.MakeList(bookmarkListDiv)

	return nil
}
