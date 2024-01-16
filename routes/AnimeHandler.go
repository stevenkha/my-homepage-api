package routes

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

type Payload struct {
	Animes []AnimeInfo `json:"data"`
}

func AnimeHandler(c *gin.Context) error {
	// TODO: Initializing cookies and client is redundant for both Manga and Anime
	// Clean it up later

	log.Info("Loading env...")

	cookieName, cookieValue, err := utils.GetEnvValues("animeCookieName", "animeCookieValue")
	if err != nil {
		return err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Debug("Could not initialize cookiejar: ")
		return err
	}

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", utils.AnimeUrl, nil)
	if err != nil {
		log.Debug("Could not create request: ")
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug("Error reading response body: ")
		return err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Debug("Error parsing html: ")
		return err
	}

	watchingListDiv := utils.GetListDiv(doc, utils.AnimeListClass)
	if watchingListDiv == nil {
		log.Debug("Could not get list of series")
	}

	animes := utils.MakeList(watchingListDiv)

	payload := formatResp(animes)

	c.JSON(http.StatusOK, payload)

	return nil
}

func formatResp(series []*html.Node) Payload {
	var anime AnimeInfo
	var resPayload Payload

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

func checkProgress(progress string) bool {
	parts := strings.Split(progress, "/")

	first, err1 := strconv.Atoi(parts[0])
	second, err2 := strconv.Atoi(parts[1])

	if err1 == nil && err2 == nil {
		return first == second
	}

	return false
}
