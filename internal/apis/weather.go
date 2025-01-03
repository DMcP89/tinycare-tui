package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var weather_url = "https://api.openweathermap.org/data/2.5/weather?zip=%s&APPID=%s&units=imperial"

// GetWeather will retrieve the current weather for a given postal code via OpenWeatherMap using its API
// https://api.openweathermap.org/data/2.5/weather?q=07070&APPID=[API KEY]
func GetWeather(postal_code string) (string, error) {
	if api_key, ok := os.LookupEnv("OPEN_WEATHER_MAP_API_KEY"); ok {
		weather := ""
		//		weather_url := "https://api.openweathermap.org/data/2.5/weather?zip=" + postal_code + "&APPID=" + api_key + "&units=imperial"
		resp, err := http.Get(fmt.Sprintf(weather_url, postal_code, api_key))
		if err != nil {
			return weather, fmt.Errorf("unable to retrieve weather data from openweathermap.org: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("unexpected response status code from OpenWeather API: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return weather, fmt.Errorf("error reading response from openweathermap.org: %w", err)
		}
		var weather_data map[string]interface{}
		err = json.Unmarshal(body, &weather_data)
		if err != nil {
			return weather, fmt.Errorf("error parsing response data from openweathermap.org: %w", err)
		}
		weather = fmt.Sprintf("%s (%s)\n", weather_data["name"], weather_data["sys"].(map[string]interface{})["country"])
		weather += fmt.Sprintf("Current: %.2f°\n", weather_data["main"].(map[string]interface{})["temp"])
		weather += fmt.Sprintf("High: %.2f°\n", weather_data["main"].(map[string]interface{})["temp_max"])
		weather += fmt.Sprintf("Low: %.2f°\n", weather_data["main"].(map[string]interface{})["temp_min"])
		weather += fmt.Sprintf("Humidity: %.2f%%\n", weather_data["main"].(map[string]interface{})["humidity"])
		weather += fmt.Sprintf("Wind: %.2f mph\n", weather_data["wind"].(map[string]interface{})["speed"])
		return weather, nil
	} else {
		return "", fmt.Errorf("OPEN_WEATHER_MAP_API_KEY environment variable not set")
	}
}
