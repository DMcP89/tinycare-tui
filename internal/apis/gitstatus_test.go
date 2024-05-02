package apis

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
)

const token = "TESTTOKEN"
const user = "DMcP89"
const page = 1
const lookback = 7

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
}

func Test_GetGitHubCommits(t *testing.T) {
	defer gock.Off()

	gock.New(reqUrl + "/user").
		Get("/").
		Reply(200).
		File("testdata/github_user.json")

	eventsEndpoint := fmt.Sprintf("/users/%s/events?per_page=100&page=%d", user, page)
	gock.New(reqUrl + eventsEndpoint).
		Get("/").
		Reply(200).
		File("testdata/github_events.json")
	_, err := GetGitHubCommits(token, lookback)
	if err != nil {
		t.Errorf("Error testing GetGitHubCommits: %s", err.Error())
	}
}
