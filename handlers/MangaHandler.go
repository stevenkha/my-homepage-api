package handlers

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

type MangaPayload struct {
	Mangas []MangaInfo `json:"mangas"`
}

func MangaHandler(c *gin.Context) {

	dataName, dataValue, err := utils.GetEnvValues("bookmarkDataName", "bookmarkDataValue")
	if err != nil {
		log.Error("Could not get env values")
	}

	client := &http.Client{}

	// form data to send to the manganato API when sending post request to get user bookmarks
	formData := url.Values{}
	formData.Set(dataName, dataValue)
	formData.Set("bm_source", "manganato")
	formData.Set("out_type", "html")

	req, err := http.NewRequest("POST", utils.MangaUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Error("Could not make post request to manganato api")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Could not read post response body")
	}

	var postPayload BookmarkResponse
	err = json.Unmarshal(body, &postPayload)
	if err != nil {
	}

	// parse only the bookmarks html portion of the response payload
	doc, err := html.Parse(strings.NewReader(string(postPayload.Data)))
	if err != nil {
		log.Error("Error parsing html: ")
	}

	bookmarkListDiv := utils.GetListDiv(doc, "user-bookmark-content")
	if bookmarkListDiv == nil {
		log.Error("Could not find bookmark list node")
	}

	mangas := utils.MakeList(bookmarkListDiv)

	payload := formatMangaResp(mangas)

	c.JSON(http.StatusOK, payload)
}

func formatMangaResp(mangas []*html.Node) MangaPayload {
	var manga MangaInfo
	var resPayload MangaPayload

	// similar note about getting animes
	// if DOM structure changes this will break but I'll worry about it later...
	for _, n := range mangas {
		manga.Cover = n.FirstChild.FirstChild.NextSibling.Attr[0].Val
		manga.Title = n.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.FirstChild.Data
		manga.Viewed = n.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild.Data
		manga.Current = n.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild.Data

		resPayload.Mangas = append(resPayload.Mangas, manga)
	}

	return resPayload
}
