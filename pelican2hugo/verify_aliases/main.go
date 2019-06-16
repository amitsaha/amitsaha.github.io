// code for findings links on a page from https://www.devdungeon.com/content/web-scraping-go#make_http_get_request
package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var links []string

// This will get called for each HTML element found
func processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	if exists {
		if strings.HasPrefix(href, "https://echorand.me/") {
			parsedUrl, err := url.Parse(href)
			if err != nil {
				log.Fatal(err)
			}
			if strings.HasSuffix(parsedUrl.Path, ".html") {
				links = append(links, parsedUrl.Path)
			}
		}
	}

}

func verifyLinks() {
	for _, l := range links {
        log.Printf("Verifying %s\n", l)
		_, err := http.Get(os.Args[2] + l)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func main() {
	// Make HTTP request
	response, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find all links and process them with the function
	// defined earlier
	document.Find("a").Each(processElement)
	verifyLinks()
}
