package utils

// Create a go function that will retrieve the current weather for a given postal code via web scraping.

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetWeather will retrieve the current weather for a given postal code via web scraping.
func GetWeather(postalCode string) (string, error) {
	// Make HTTP request
	reqUrl := "https://weather.com/weather/today/l/" + postalCode + ":4:US"
	response, err := http.Get(reqUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	// Get the weather
	// The weather is not contained in a div called "today_nowcard-phrase" anymore
	// the current temp is stored in a span with class CurrentConditions--tempValue--MHmYY
	weather := document.Find("span.CurrentConditions--tempValue--MHmYY").First().Text()
	location := document.Find("h1.CurrentConditions--location--1YWj_").First().Text()
	weather = strings.TrimSpace(weather)

	// Print the weather
	return fmt.Sprintf("%s %s", location, weather), nil
}
