package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GetTodaysTasks will retrieve the tasks for today by querying the Todoist API.
func GetTodaysTasks() (string, error) {
	// Make HTTP request
	reqURL := "https://api.todoist.com/rest/v2/tasks?filter=today"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	// Set the API token as a header
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TODOIST_TOKEN"))

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal the JSON data
	var tasks []map[string]interface{}
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return "", err
	}

	// Print the tasks
	var output string
	for _, task := range tasks {
		fmt.Printf("Task: %v\n", task)
		output += fmt.Sprintf("* %s\n", task["content"])
	}

	return output, nil
}
