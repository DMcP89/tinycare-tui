package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GetWeather will retrieve the current weather for a given postal code via OpenWeatherMap using its API
// https://api.openweathermap.org/data/2.5/weather?q=07070&APPID=[API KEY]
func GetWeather(postal_code string) (string, error) {
	weather := ""
	weather_url := "https://api.openweathermap.org/data/2.5/weather?zip=" + postal_code + "&APPID=" + os.Getenv("OPEN_WEATHER_MAP_API_KEY") + "&units=imperial"
	resp, err := http.Get(weather_url)
	if err != nil {
		return weather, err
	}
	defer resp.Body.Close()
	// replace ioutil.ReadAll with io.Copy
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weather, err
	}
	var weather_data map[string]interface{}
	err = json.Unmarshal(body, &weather_data)
	if err != nil {
		return weather, err
	}
	weather = fmt.Sprintf("%s (%s)\n", weather_data["name"], weather_data["sys"].(map[string]interface{})["country"])
	weather += fmt.Sprintf("Current: %.2f°\n", weather_data["main"].(map[string]interface{})["temp"])
	weather += fmt.Sprintf("High: %.2f°\n", weather_data["main"].(map[string]interface{})["temp_max"])
	weather += fmt.Sprintf("Low: %.2f°\n", weather_data["main"].(map[string]interface{})["temp_min"])
	weather += fmt.Sprintf("Humidity: %.2f%%\n", weather_data["main"].(map[string]interface{})["humidity"])
	weather += fmt.Sprintf("Wind: %.2f mph\n", weather_data["wind"].(map[string]interface{})["speed"])
	return weather, nil
}
