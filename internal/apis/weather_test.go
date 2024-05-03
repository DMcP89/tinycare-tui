package apis

import (
	"fmt"
	"os"
	"testing"

	"github.com/h2non/gock"
)

var api_key string

func init() {
	api_key, ok := os.LookupEnv("OPEN_WEATHER_MAP_API_KEY")
	if !ok {
		api_key = "TESTAPIKEY"
		os.Setenv("OPEN_WEATHER_MAP_API_KEY", api_key)
	}
}

func TestValidZipCode(t *testing.T) {
	defer gock.Off()
	postal_code := "10005"
	gock.New(fmt.Sprintf(weather_url, postal_code, api_key)).
		Get("/").
		Reply(200).
		//		JSON(map[string]string{"foo": "bar"})
		File("testdata/weather.json")

	_, err := GetWeather(postal_code)
	if err != nil {
		t.Error(err)
	}
}

func TestEmptyStringZipCode(t *testing.T) {
	_, err := GetWeather("")
	if err == nil {
		t.Fatalf("Expected error for empty string zip")
	}
}

func TestInvalidZipCode(t *testing.T) {
	_, err := GetWeather("ABCDEF")
	if err == nil {
		t.Fatalf("Expected error for invalid zip")
	}
}
