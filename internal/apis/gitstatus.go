/*
Prompt that was the basis for gitstatus.go used with ChatGPT

As a Golang Developer, create a script that will return information commits in git repositories
This script should contain 2 functions: GetDailyCommits and GetWeeklyCommits
The GetDailyCommits function should take in a path and return all of the commits for the git repositories it finds under that path from the current day
the return value from GetDailyCommits should be as follows:

	RepoName1
	Commit Hash - Commit Message (# of hours ago the commit was made)

	RepoName2
	Commit Hash - Commit Message (# of hours ago the commit was made)
	Commit Hash - Commit Message (# of hours ago the commit was made)

For example the return for GetDailyCommits when passed a path with 1 git repository under it named Harambot:

	Harambot
	1655c81 - Added new main.tf for azure deployments (1 hour ago)
	5bc9bdd - Fixed an issue with parsing team names (12 hours ago)

The GetWeeklyCommits should be similar however instead of only returning commits from the last day it should return all the commits from the past 7 days. Its output should be formated the same way GetDailyCommits is
*/
package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DMcP89/tinycare-tui/internal/utils"
)

// Actor represents the actor in the JSON structure.
type Actor struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarID   string `json:"gravatar_id"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatar_url"`
}

// Repo represents the repo in the JSON structure.
type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Author represents the author in the Commit structure.
type Author struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Commit represents a commit in the JSON structure.
type Commit struct {
	SHA      string `json:"sha"`
	Author   Author `json:"author"`
	Message  string `json:"message"`
	Distinct bool   `json:"distinct"`
	URL      string `json:"url"`
}

// Payload represents the payload in the JSON structure.
type Payload struct {
	RepositoryID int      `json:"repository_id"`
	PushID       int64    `json:"push_id"`
	Size         int      `json:"size"`
	DistinctSize int      `json:"distinct_size"`
	Ref          string   `json:"ref"`
	Head         string   `json:"head"`
	Before       string   `json:"before"`
	Commits      []Commit `json:"commits"`
}

// Event represents the main structure of each event in the JSON array.
type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
}

const reqUrl = "https://api.github.com"
const missingTokenMessage = "GITHUB_TOKEN environment variable not set correctly"

func GetGitHubUser(token string) (string, error) {
	userEndpoint := "/user"
	// get the username of the authenticated user
	req, err := http.NewRequest("GET", reqUrl+userEndpoint, nil)

	if err != nil {
		return "", err
	}
	// Set the API token as a header
	req.Header.Set("Authorization", "Bearer "+token)

	body, err := utils.SendRequest(req)
	if err != nil {
		return "", err
	}

	var user map[string]interface{}
	err = json.Unmarshal(body, &user)
	return user["login"].(string), err
}

func GetGitHubEvents(token string, login string, page int) ([]Event, error) {
	eventsEndpoint := fmt.Sprintf("/users/%s/events?per_page=100&page=%d", login, page)

	req, err := http.NewRequest("GET", reqUrl+eventsEndpoint, nil)

	if err != nil {
		return nil, err
	}
	// Set the API token as a header
	req.Header.Set("Authorization", "Bearer "+token)

	body, err := utils.SendRequest(req)
	if err != nil {
		return nil, err
	}
	// Unmarshal the JSON data

	var events []Event
	err = json.Unmarshal(body, &events)

	return events, err
}

func GetGitHubCommits(token string) (string, string, error) {

	if token != "" {
		user, userErr := GetGitHubUser(token)
		if userErr != nil {
			return "", "", fmt.Errorf("unable to get Github User: %w", userErr)
		}
		var totalEvents []Event

		dayLookBackTime := time.Now().AddDate(0, 0, -1)
		weekLookBackTime := time.Now().AddDate(0, 0, -7)
		page := 1
		for {
			events, eventsErr := GetGitHubEvents(token, user, page)

			if eventsErr != nil {
				return "", "", fmt.Errorf("unable to get events for user %s: %w", user, eventsErr)
			}
			totalEvents = append(totalEvents, events...)
			page++
			if events[len(events)-1].CreatedAt.Before(weekLookBackTime) {
				break
			}
		}

		var weekOutput string
		var dayOutput string

		pullCommits := func(event Event) string {
			var result string
			for _, commit := range event.Payload.Commits {
				timeSinceCommit := time.Since(event.CreatedAt.In(time.Local))
				formattedTimeSinceCommit := utils.HumanizeDuration(timeSinceCommit)
				result += fmt.Sprintf("[yello]%s[white] (%s)\n", commit.Message, formattedTimeSinceCommit)
			}
			return result
		}

		for _, event := range totalEvents {
			if len(event.Payload.Commits) > 0 && event.CreatedAt.In(time.Local).After(dayLookBackTime) {
				commitText := pullCommits(event)
				dayOutput += fmt.Sprintf("[red]%s[white]\n", event.Repo.Name)
				dayOutput += commitText
				weekOutput += fmt.Sprintf("[red]%s[white]\n", event.Repo.Name)
				weekOutput += commitText
			} else if len(event.Payload.Commits) > 0 && event.CreatedAt.In(time.Local).After(weekLookBackTime) {
				weekOutput += fmt.Sprintf("[red]%s[white]\n", event.Repo.Name)
				weekOutput += pullCommits(event)
			}

		}

		return dayOutput, weekOutput, nil

	} else {
		return missingTokenMessage, missingTokenMessage, nil
	}

}
