package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"

	"my-homepage-api/utils"
)

type ScheduledAnimesResp struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
}

type AnimeInfo struct {
	Cover   string `json:"cover"`
	Title   string `json:"title"`
	Viewed  string `json:"viewed"`
	Current string `json:"current"`
	Slug    string `json:"slug"`
}

type AnimePayload struct {
	ScheduledAnimes []AnimeInfo `json:"scheduledAnimes"`
	BacklogAnimes   []AnimeInfo `json:"backlogAnimes"`
}

func AnimeHandler(c *gin.Context) {
	client := &http.Client{}

	scheduledAnimes := getScheduledAnimes(client)
	bookmarkedAnimes := getBookmarkedAnimes(client)

	var anime AnimeInfo
	var resPayload AnimePayload

	for _, n := range bookmarkedAnimes {
		anime.Cover = n.FirstChild.FirstChild.FirstChild.Attr[0].Val
		formatCover(&anime.Cover)
		anime.Title = strings.TrimSpace(n.FirstChild.NextSibling.FirstChild.FirstChild.Data)
		progress := n.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.Data

		slug := n.FirstChild.NextSibling.FirstChild.Attr[0].Val
		parts := strings.Split(slug, "/")
		anime.Slug = parts[2]

		if newEpisode(scheduledAnimes, anime.Title) {
			resPayload.ScheduledAnimes = append(resPayload.ScheduledAnimes, anime)
		}

		parts = strings.Split(progress, "/")
		anime.Viewed = parts[0]
		anime.Current = parts[1]

		caughtUp, err := checkProgress(anime.Current, anime.Viewed)
		if err != nil {
			log.Error("Could not check progress: %v", err)
		}

		if !caughtUp {
			resPayload.BacklogAnimes = append(resPayload.BacklogAnimes, anime)
		}
	}

	c.JSON(http.StatusOK, resPayload)
}

func UpdateHandler(c *gin.Context) {
	cookieName, cookieValue, err := utils.GetEnvValues("animeCookieName", "animeCookieValue")
	if err != nil {
		log.Error("Could not get env values")
	}

	c.SetCookie(cookieName, cookieValue, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, utils.BookmarkedAnimesUrl)
}

func newEpisode(scheduledAnimes []string, title string) bool {
	for _, a := range scheduledAnimes {
		if strings.Contains(a, title) {
			return true
		}
	}

	return false
}

func getBookmarkedAnimes(client *http.Client) []*html.Node {
	cookieName, cookieValue, err := utils.GetEnvValues("animeCookieName", "animeCookieValue")
	if err != nil {
		log.Error("Could not get env values")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error("Could not initialize cookiejar: ")
	}

	client.Jar = jar

	req, err := http.NewRequest("GET", utils.BookmarkedAnimesUrl, nil)
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

	return utils.MakeList(watchingListDiv)
}

func getScheduledAnimes(client *http.Client) []string {
	req, err := http.NewRequest("GET", utils.ScheduledAnimeUrl, nil)
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

	var schAnimeResp []ScheduledAnimesResp
	err = json.Unmarshal(body, &schAnimeResp)
	if err != nil {
		log.Error("Could not unmarshal scheduled animes response")
	}

	dateFormat := "2006-01-02"
	today := time.Now().Format(dateFormat)
	schAnimes := make([]string, 0)

	// only add animes that released on current day
	for _, a := range schAnimeResp {
		animeRelease := time.Unix(a.Timestamp, 0).Format(dateFormat)
		if strings.Contains(today, animeRelease) {
			schAnimes = append(schAnimes, a.Title)
		}
	}

	return schAnimes
}

func formatAnimeResp(series []*html.Node) AnimePayload {
	var anime AnimeInfo
	var resPayload AnimePayload

	for _, n := range series {
		anime.Title = n.FirstChild.Data
	}

	// scuffed way of getting the data from the html tags
	// if DOM structure changes this will break but I'll worry about it later...

	// for _, n := range series {
	// 	anime.Cover = n.FirstChild.FirstChild.FirstChild.Attr[0].Val
	// 	anime.Title = n.FirstChild.NextSibling.FirstChild.FirstChild.Data

	// 	formatCover(&anime.Cover)
	// 	resPayload.Animes = append(resPayload.Animes, anime)
	// }

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
