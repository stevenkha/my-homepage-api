package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"

	"my-homepage-api/utils"
)

type AnimeInfo struct {
	Cover string `json:"cover"`
	Title string `json:"title"`
}

type AnimePayload struct {
	Animes []AnimeInfo `json:"animes"`
}

func AnimeHandler(c *gin.Context) {

	client := &http.Client{}

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Error("Error getting location")
	}

	updatedUrl := utils.AnimeUrl + time.Now().In(loc).Format("2006-01-02")

	req, err := http.NewRequest("GET", updatedUrl, nil)
	if err != nil {
		log.Error("Could not create request: ")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Could not send HTTP request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body: ")
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Error("Error parsing html: ")
	}

	c.JSON(http.StatusOK, payload)
}

func formatAnimeResp(series []*html.Node) AnimePayload {
	var anime AnimeInfo
	var resPayload AnimePayload

	// scuffed way of getting the data from the html tags
	// if DOM structure changes this will break but I'll worry about it later...
	for _, n := range series {
		anime.Cover = n.FirstChild.FirstChild.FirstChild.Attr[0].Val
		anime.Title = n.FirstChild.NextSibling.FirstChild.FirstChild.Data

		formatCover(&anime.Cover)
		resPayload.Animes = append(resPayload.Animes, anime)
	}

	return resPayload
}

/*
Format the cover image of the animes to be higher resolution
*/
func formatCover(cover *string) {
	w := 225
	h := 300

	// Split the URL into urlParts based on '/'
	urlParts := strings.Split(*cover, "/")

	// get the width and height string from the url array
	dimensionParts := strings.Split(urlParts[4], ",")

	// set new width and height
	dimensionParts[0] = fmt.Sprintf("w_%d", w)
	dimensionParts[1] = fmt.Sprintf("h_%d", h)

	// combine
	urlParts[4] = strings.Join(dimensionParts, ",")

	*cover = strings.Join(urlParts, "/")
}

/*
Check if I am caught up with the latest episode
format of this data is '12/12' so check if the first half is equal to the second i.e '11/12'
*/
func checkProgress(current string, viewed string) (bool, error) {
	c, err1 := strconv.Atoi(current)
	if err1 != nil {
		return false, err1
	}

	v, err2 := strconv.Atoi(viewed)
	if err2 != nil {
		return false, err2
	}

	return c == v, nil
}
