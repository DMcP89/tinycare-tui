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
	"context"
	"fmt"
	"time"

	"github.com/DMcP89/tinycare-tui/internal/utils"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

const missingTokenMessage = "GITHUB_TOKEN environment variable not set correctly"

// newGitHubClient creates a new GitHub client with authentication
func newGitHubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetGitHubUser(token string) (string, error) {
	client := newGitHubClient(token)
	ctx := context.Background()
	
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", err
	}
	
	if user.Login == nil {
		return "", fmt.Errorf("user login not found")
	}
	
	return *user.Login, nil
}

func GetGitHubEvents(token string, login string, page int) ([]*github.Event, error) {
	client := newGitHubClient(token)
	ctx := context.Background()
	
	opts := &github.ListOptions{
		Page:    page,
		PerPage: 100,
	}
	
	events, _, err := client.Activity.ListEventsPerformedByUser(ctx, login, false, opts)
	if err != nil {
		return nil, err
	}
	
	return events, nil
}

func GetGitHubCommits(token string) (string, string, error) {

	if token != "" {
		user, userErr := GetGitHubUser(token)
		if userErr != nil {
			return "", "", fmt.Errorf("unable to get Github User: %w", userErr)
		}
		var totalEvents []*github.Event

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
			if len(events) > 0 && events[len(events)-1].CreatedAt != nil && events[len(events)-1].CreatedAt.Before(weekLookBackTime) {
				break
			}
		}

		var weekOutput string
		var dayOutput string

		pullCommits := func(event *github.Event) string {
			var result string
			// Type assertion to get push event payload
			if event.Type != nil && *event.Type == "PushEvent" {
				// Parse the raw payload as PushEvent
				if pushEvent, ok := event.Payload().(*github.PushEvent); ok && pushEvent.Commits != nil {
					for _, commit := range pushEvent.Commits {
						if event.CreatedAt != nil && commit.Message != nil {
							timeSinceCommit := time.Since(event.CreatedAt.In(time.Local))
							formattedTimeSinceCommit := utils.HumanizeDuration(timeSinceCommit)
							result += fmt.Sprintf("[yello]%s[white] (%s)\n", *commit.Message, formattedTimeSinceCommit)
						}
					}
				}
			}
			return result
		}

		for _, event := range totalEvents {
			if event.Type != nil && *event.Type == "PushEvent" && event.CreatedAt != nil {
				// Check if this is a push event with commits and within our time range
				if pushEvent, ok := event.Payload().(*github.PushEvent); ok && pushEvent.Commits != nil && len(pushEvent.Commits) > 0 {
					if event.CreatedAt.In(time.Local).After(dayLookBackTime) {
						commitText := pullCommits(event)
						if event.Repo != nil && event.Repo.Name != nil {
							dayOutput += fmt.Sprintf("[red]%s[white]\n", *event.Repo.Name)
							dayOutput += commitText
							weekOutput += fmt.Sprintf("[red]%s[white]\n", *event.Repo.Name)
							weekOutput += commitText
						}
					} else if event.CreatedAt.In(time.Local).After(weekLookBackTime) {
						if event.Repo != nil && event.Repo.Name != nil {
							weekOutput += fmt.Sprintf("[red]%s[white]\n", *event.Repo.Name)
							weekOutput += pullCommits(event)
						}
					}
				}
			}
		}

		return dayOutput, weekOutput, nil

	} else {
		return missingTokenMessage, missingTokenMessage, nil
	}

}
