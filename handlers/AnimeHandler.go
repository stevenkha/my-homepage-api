package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"

	"my-homepage-api/utils"
)

type AnimeInfo struct {
	Cover  string `json:"cover"`
	Title  string `json:"title"`
	Latest bool   `json:"latest"`
}

type AnimePayload struct {
	Animes []AnimeInfo `json:"animes"`
}

func AnimeHandler(c *gin.Context) {
	cookieName, cookieValue, err := utils.GetEnvValues("animeCookieName", "animeCookieValue")
	if err != nil {
		log.Error("Could not get env values")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error("Could not initialize cookiejar: ")
	}

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", utils.AnimeUrl, nil)
	if err != nil {
		log.Error("Could not create request: ")
	}

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: cookieValue,
	}

	req.AddCookie(cookie)

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

	watchingListDiv := utils.GetListDiv(doc, utils.AnimeListClass)
	if watchingListDiv == nil {
		log.Error("Could not get list of series")
	}

	animes := utils.MakeList(watchingListDiv)

	payload := formatAnimeResp(animes)

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
		anime.Latest = checkProgress(n.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.Data)

		resPayload.Animes = append(resPayload.Animes, anime)
	}

	return resPayload
}

// Check if I am caught up with the latest episode
// format of this data is '12/12' so check if the first half is equal to the second i.e '11/12'
func checkProgress(progress string) bool {
	parts := strings.Split(progress, "/")

	first, err1 := strconv.Atoi(parts[0])
	second, err2 := strconv.Atoi(parts[1])

	if err1 == nil && err2 == nil {
		return first == second
	}

	return false
}
