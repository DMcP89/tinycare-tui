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
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

func GetGitHubCommits() (string, error) {
	if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		reqUrl := "https://api.github.com"
		userEndpoint := "/user"
		// get the username of the authenticated user
		req, err := http.NewRequest("GET", reqUrl+userEndpoint, nil)

		if err != nil {
			return "Error Creating request", err
		}
		// Set the API token as a header
		req.Header.Set("Authorization", "Bearer "+token)

		// Send HTTP request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "Error sending request", err
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			return "Non 200 response", fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Error reading the body", err
		}

		// Unmarshal the JSON data

		var user map[string]interface{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			return "Error unmarshalling user", err
		}

		// get the events of the authenticated user
		eventsEndpoint := fmt.Sprintf("/users/%s/events?per_page=2", user["login"])

		req, err = http.NewRequest("GET", reqUrl+eventsEndpoint, nil)

		if err != nil {
			return "Error Creating request", err
		}
		// Set the API token as a header
		req.Header.Set("Authorization", "Bearer "+token)

		// Send HTTP request
		resp, err = client.Do(req)
		if err != nil {
			return "Error sending request", err
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			return "Non 200 response", fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
		}

		// Read the response body
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return "Error reading the body", err
		}

		// Unmarshal the JSON data

		var events []Event
		err = json.Unmarshal(body, &events)
		if err != nil {
			return "Error unmarshalling user", err
		}

		// Print the events
		var output string
		for _, event := range events {
			for _, commit := range event.Payload.Commits {
				output += fmt.Sprintf("%s - %s (%d)", event.Repo.Name, commit.Message, event.CreatedAt)
			}
		}
		return output, nil

	} else {
		return "No Token set", nil
	}
	return "", nil
}

func GetDailyCommits(path string) (string, error) {
	repositories, err := findGitRepositories(path)
	if err != nil {
		return "", err
	}

	result := ""
	for _, repo := range repositories {
		commitMessages, err := getCommitsFromTimeRange(repo, time.Now().AddDate(0, 0, -1), time.Now())
		if err != nil {
			return "", err
		}

		if len(commitMessages) > 0 {
			result += fmt.Sprintf("[red]%s[white]\n", repo)
			result += commitMessages + "\n"
		}
	}

	return result, nil
}

func GetWeeklyCommits(path string) (string, error) {
	repositories, err := findGitRepositories(path)
	if err != nil {
		return "", err
	}

	result := ""
	for _, repo := range repositories {
		commitMessages, err := getCommitsFromTimeRange(repo, time.Now().AddDate(0, 0, -7), time.Now())
		if err != nil {
			return "", err
		}

		if len(commitMessages) > 0 {
			result += fmt.Sprintf("[red]%s[white]\n", repo)
			result += commitMessages + "\n"
		}
	}

	return result, nil
}

func findGitRepositories(path string) ([]string, error) {
	var repositories []string
	//split the path into a slice of strings by comma

	paths := strings.Split(path, ",")
	for _, path := range paths {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".git" {
				repositories = append(repositories, filepath.Dir(p))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return repositories, nil
}

func getCommitsFromTimeRange(repoPath string, since time.Time, until time.Time) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	commitMessages := ""

	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}

	commitIter, err := repo.Log(&git.LogOptions{
		From: headRef.Hash(),
	})
	if err != nil {
		return "", err
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Committer.When.After(since) && commit.Committer.When.Before(until) {
			timeSinceCommit := time.Since(commit.Committer.When)
			formattedTimeSinceCommit := humanizeDuration(timeSinceCommit)

			commitMessages += fmt.Sprintf("[yellow]%s[white] - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return commitMessages, nil
}

func humanizeDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	return fmt.Sprintf("[green]%d h(s) ago[white]", hours)
}
