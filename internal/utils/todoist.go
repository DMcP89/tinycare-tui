package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetTasks() (string, error) {
	// Check for the existance of the environment variable TODOIST_TOKEN
	if token, ok := os.LookupEnv("TODOIST_TOKEN"); ok {
		return GetTodaysTasks(token)
	} else {
		return GetLocalTasks()
	}
}

func GetLocalTasks() (string, error) {
	// Check for the existance of the environment variable TODO_FILE
	// If it exists return its contents as a string
	// If it does not exist return "Please set your TODO_FILE variable"
	if todoFile, ok := os.LookupEnv("TODO_FILE"); ok {
		file, err := os.Open(todoFile)
		if err != nil {
			return "", fmt.Errorf("Unable to open %s : %w", todoFile, err)
		}
		defer file.Close()

		var output string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			output += fmt.Sprintf("☐ %s\n", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("Unable to read %s : %w", todoFile, err)
		}
		return output, nil
	} else {
		return "", fmt.Errorf("No TODO_FILE environment variable set")
	}
}

// GetTodaysTasks will retrieve the tasks for today by querying the Todoist API.
func GetTodaysTasks(token string) (string, error) {
	// Make HTTP request
	reqURL := "https://api.todoist.com/rest/v2/tasks?filter=today"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("Unable to create request for Todoist: %w", err)
	}

	// Set the API token as a header
	req.Header.Set("Authorization", "Bearer "+token)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error while sending request to Todoist: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status code from Todoist: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to read response data: %w", err)
	}

	// Unmarshal the JSON data
	var tasks []map[string]interface{}
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal response data: %w", err)
	}

	// Print the tasks
	var output string
	for _, task := range tasks {
		output += fmt.Sprintf("☐ %s\n", task["content"])
	}
	completed, err := GetCompletedTasks(token)
	if err != nil {
		output += "Could not fetch completed tasks \n"
		output += err.Error()
		output += completed
		return output, nil
	}
	output += completed
	return output, nil
}

func GetCompletedTasks(token string) (string, error) {
	reqURL := "https://api.todoist.com/sync/v9/completed/get_all"
	today := strings.Split(time.Now().Format(time.RFC3339), "T")[0]

	var output string
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request for completed tasks: %w", err)
	}
	q := req.URL.Query()
	q.Add("since", today+"T00:00:00")
	req.URL.RawQuery = q.Encode()

	// Set the API token as a header
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request to Todoist for completed tasks: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading repsonse data from Todoist for completed tasks: %w", err)
	}

	// Unmarshal the JSON data
	var tasks map[string]interface{}
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling response data from Todoist for completed tasks: %w", err)
	}
	// Print the tasks

	if items, ok := tasks["items"].([]interface{}); ok {
		for _, task := range items {
			task, ok := task.(map[string]interface{})
			if !ok {
				continue
			}
			//fmt.Println(task["content"])
			output += fmt.Sprintf("✅ %s\n", task["content"])
		}
	}
	return output, nil

}
