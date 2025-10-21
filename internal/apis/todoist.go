package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DMcP89/tinycare-tui/internal/utils"
)

// GetTodaysTasks will retrieve the tasks for today by querying the Todoist API.
func GetTodaysTasks(token string) (string, error) {
	// Make HTTP request
	reqURL := "https://api.todoist.com/rest/v2/tasks?filter=today|overdue"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("unable to create request for Todoist: %w", err)
	}

	// Set the API token as a header
	req.Header.Set("Authorization", "Bearer "+token)

	// Read the response body
	body, err := utils.SendRequest(req)
	if err != nil {
		return "", fmt.Errorf("unable to read response data: %w", err)
	}

	// Unmarshal the JSON data
	var tasks []map[string]interface{}
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal response data: %w", err)
	}

	// Print the tasks
	var output string
	today := strings.Split(time.Now().Format(time.RFC3339), "T")[0]
	for _, task := range tasks {
		content := task["content"]
		// Check if task is overdue
		isOverdue := false
		if due, ok := task["due"].(map[string]interface{}); ok {
			if dueDate, ok := due["date"].(string); ok {
				// Compare due date with today
				if dueDate < today {
					isOverdue = true
				}
			}
		}

		if isOverdue {
			output += fmt.Sprintf("[red]☐ %s[white]\n", content)
		} else {
			output += fmt.Sprintf("☐ %s\n", content)
		}
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
		return "", fmt.Errorf("error creating request for completed tasks: %w", err)
	}
	q := req.URL.Query()
	q.Add("since", today+"T00:00:00")
	req.URL.RawQuery = q.Encode()

	// Set the API token as a header
	req.Header.Set("Authorization", "Bearer "+token)

	// Read the response body
	body, err := utils.SendRequest(req)
	if err != nil {
		return "", fmt.Errorf("error reading response data from Todoist for completed tasks: %w", err)
	}

	// Unmarshal the JSON data
	var tasks map[string]interface{}
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response data from Todoist for completed tasks: %w", err)
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
