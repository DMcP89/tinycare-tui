package apis

import (
	"net/http"
	"log/slog"

	"github.com/DMcP89/tinycare-tui/internal/utils"
)

const jokeUrl = "https://icanhazdadjoke.com/"

func GetJoke() (string, error) {
	req, err := http.NewRequest("GET", jokeUrl, nil)

	if err != nil {
		return "", err
	}
	// Set the API token as a header
	req.Header.Set("Accept", "text/plain")

	joke, err := utils.SendRequest(req)
	if err != nil {
		slog.Error("failed to fetch joke from server", "error", err)
		return "", err
	}

	return string(joke[:]), nil
}
