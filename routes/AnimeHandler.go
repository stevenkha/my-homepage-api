package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	log "github.com/ccpaging/log4go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

const (
	url       = "https://yugenanime.tv/mylist/"
	listClass = "list-entries"
)

type payload struct {
	Title  string `json:"title"`
	Latest bool   `json:"latest"`
	Image  string `json:"image"`
}

func AnimeHandler(c *gin.Context) error {

	log.Info("Loading env...")

	cookieName, cookieValue, err := getEnvValues()
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

	req, err := http.NewRequest("GET", url, nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug("Error reading response body:", err)
		return err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Debug("Error parsing html: ")
		return err
	}

	seriesList := getList(doc, listClass)
	if seriesList == nil {
		log.Error("Could not get list of series")
	}

	return nil
}

func getList(n *html.Node, targetClass string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, targetClass) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := getList(c, targetClass); result != nil {
			return result
		}
	}

	return nil
}

func getEnvValues() (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Debug("Error loading .env file: %s", err)
		return "", "", err
	}

	cookieName := os.Getenv("cookieName")
	if cookieName == "" {
		log.Debug("Could not read cookie name")
		return "", "", err
	}

	cookieValue := os.Getenv("cookieValue")
	if cookieValue == "" {
		log.Fatal("Could not read cookie value")
	}

	return cookieName, cookieValue, nil
}
