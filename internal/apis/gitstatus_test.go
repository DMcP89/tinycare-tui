package apis

import (
	"fmt"
	"testing"

	"github.com/h2non/gock"
)

const token = "TESTTOKEN"
const user = "DMcP89"
const page = 1
const apiUrl = "https://api.github.com"

func Test_GetGitHubUser(t *testing.T) {
	defer gock.Off()

	tests := []struct {
		name         string
		mockReply    func()
		expectedUser string
		expectError  bool
	}{
		{
			name: "Valid User",
			mockReply: func() {
				gock.New(apiUrl).
					Get("/user").
					Reply(200).
					File("testdata/github_user.json")
			},
			expectedUser: user,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			login, err := GetGitHubUser(token)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if login != tt.expectedUser {
				t.Errorf("GetGitHubUser = %s, expected %s", login, tt.expectedUser)
			}
		})
	}
}

func Test_GetGitHubEvents(t *testing.T) {
	defer gock.Off()

	tests := []struct {
		name        string
		mockReply   func()
		expectError bool
	}{
		{
			name: "Valid Events",
			mockReply: func() {
				eventsEndpoint := fmt.Sprintf("/users/%s/events", user)
				gock.New(apiUrl).
					Get(eventsEndpoint).
					MatchParam("per_page", "100").
					MatchParam("page", fmt.Sprintf("%d", page)).
					Reply(200).
					File("testdata/github_events.json")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			_, err := GetGitHubEvents(token, user, page)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func Test_GetGitHubCommits(t *testing.T) {
	defer gock.Off()

	tests := []struct {
		name        string
		mockReplies []func()
		expectError bool
	}{
		{
			name: "Valid Commits",
			mockReplies: []func(){
				func() {
					gock.New(apiUrl).
						Get("/user").
						Reply(200).
						File("testdata/github_user.json")
				},
				func() {
					eventsEndpoint := fmt.Sprintf("/users/%s/events", user)
					gock.New(apiUrl).
						Get(eventsEndpoint).
						MatchParam("per_page", "100").
						MatchParam("page", fmt.Sprintf("%d", page)).
						Reply(200).
						File("testdata/github_events.json")
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mockReply := range tt.mockReplies {
				mockReply()
			}
			_, _, err := GetGitHubCommits(token)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func Test_GetGitHubCommits_EmptyToken(t *testing.T) {
	day, week, err := GetGitHubCommits("")

	expectedMessage := "GITHUB_TOKEN environment variable not set correctly"

	if day != expectedMessage {
		t.Errorf("Expected day output: %s, got: %s", expectedMessage, day)
	}

	if week != expectedMessage {
		t.Errorf("Expected week output: %s, got: %s", expectedMessage, week)
	}

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
