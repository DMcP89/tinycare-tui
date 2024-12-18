package apis

import (
	"testing"

	"github.com/h2non/gock"
)

func TestGetTodaysTasks(t *testing.T) {
	defer gock.Off()

	tests := []struct {
		name        string
		mockReply   func()
		expectError bool
	}{
		{
			name: "Valid Today's Tasks",
			mockReply: func() {
				reqURL := "https://api.todoist.com/rest/v2/tasks?filter=today"
				gock.New(reqURL).
					Get("/").
					Reply(200).
					File("testdata/today_task.json")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			_, err := GetTodaysTasks("test-token")
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGetCompletedTasks(t *testing.T) {
	defer gock.Off()

	tests := []struct {
		name        string
		mockReply   func()
		expectError bool
	}{
		{
			name: "Valid Completed Tasks",
			mockReply: func() {
				reqURL := "https://api.todoist.com/sync/v9/completed/get_all"
				gock.New(reqURL).
					Get("/").
					Reply(200).
					File("testdata/complete_task.json")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			_, err := GetCompletedTasks("test-token")
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
