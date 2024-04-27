package utils

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
)

const token = "TESTTOKEN"
const user = "DMcP89"
const page = 1

func Test_GetGitHubUser(t *testing.T) {
	defer gock.Off()

	gock.New(reqUrl + "/user").
		Get("/").
		Reply(200).
		File("testdata/github_user.json")

	login, err := GetGitHubUser(token)
	if err != nil {
		t.Errorf("Error testing GetGitHubUser")
	}

	if login != user {
		t.Errorf("GetGitHubUser = %s, epexcted DMcP89", login)
	}
}

func Test_GetGitHubEvents(t *testing.T) {
	defer gock.Off()

	eventsEndpoint := fmt.Sprintf("/users/%s/events?per_page=100&page=%d", user, page)
	gock.New(reqUrl + eventsEndpoint).
		Get("/").
		Reply(200).
		File("testdata/github_events.json")

	_, err := GetGitHubEvents(token, user, page)
	if err != nil {
		t.Errorf("Error testing GetGitHubEvents: %s", err.Error())
	}

	t.Fatalf("Test not implemented")
}

func Test_GetGitHubCommits(t *testing.T) {
	t.Fatalf("Test not implemented")
}

func Test_GetDailyCommits(t *testing.T) {
	t.Fatalf("Test not implemented")
}

func Test_GetWeeklyCommits(t *testing.T) {
	t.Fatalf("Test not implemented")
}

func Test_GetRepos(t *testing.T) {
	t.Fatalf("Test not implemented")
}

func Test_FindGitRepositories(t *testing.T) {
	t.Fatalf("Test not implemented")
}

func Test_GetCommitsFromTimeRange(t *testing.T) {
	t.Fatalf("Test not implemented")
}
func Test_HumanizeDuration(t *testing.T) {
	t.Fatalf("Test not implemented")
}
