package utils

import (
	"os"
	"strings"

	log "github.com/ccpaging/log4go"
	"golang.org/x/net/html"

	"github.com/joho/godotenv"
)

func GetEnvValues(name string, value string) (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Debug("Error loading .env file: %s", err)
		return "", "", err
	}

	cookieName := os.Getenv(name)
	if cookieName == "" {
		log.Debug("Could not read cookie name")
		return "", "", err
	}

	cookieValue := os.Getenv(value)
	if cookieValue == "" {
		log.Fatal("Could not read cookie value")
	}

	return cookieName, cookieValue, nil
}

func GetListDiv(n *html.Node, targetClass string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, targetClass) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := GetListDiv(c, targetClass); result != nil {
			return result
		}
	}

	return nil
}
