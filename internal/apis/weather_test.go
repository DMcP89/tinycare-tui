package apis

import (
	"fmt"
	"os"
	"testing"

	"github.com/h2non/gock"
)

var api_key string

func init() {
	_, ok := os.LookupEnv("OPEN_WEATHER_MAP_API_KEY")
	if !ok {
		os.Setenv("OPEN_WEATHER_MAP_API_KEY", "TESTAPIKEY")
	}
}

func TestGetWeather(t *testing.T) {
	defer gock.Off()

	tests := []struct {
		name        string
		postal_code string
		mockReply   func()
		expectError bool
	}{
		{
			name:        "Valid Zip Code",
			postal_code: "10005",
			mockReply: func() {
				gock.New(fmt.Sprintf(weather_url, "10005", api_key)).
					Get("/").
					Reply(200).
					File("testdata/weather.json")
			},
			expectError: false,
		},
		{
			name:        "Empty String Zip Code",
			postal_code: "",
			mockReply:   func() {},
			expectError: true,
		},
		{
			name:        "Invalid Zip Code",
			postal_code: "ABCDEF",
			mockReply:   func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			_, err := GetWeather(tt.postal_code)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
