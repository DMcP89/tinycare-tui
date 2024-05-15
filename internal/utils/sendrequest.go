package utils

import (
	"fmt"
	"io"
	"net/http"
)

func SendRequest(req *http.Request) ([]byte, error) {
	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	// Read the response body
	return io.ReadAll(resp.Body)
}
