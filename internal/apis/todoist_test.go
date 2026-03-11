package apis

import (
	"strings"
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
				gock.New("https://api.todoist.com").
					Get("/api/v1/tasks/filter").
					MatchParam("query", "today|overdue").
					Reply(200).
					File("testdata/today_task.json")
				
				gock.New("https://api.todoist.com").
					Get("/api/v1/tasks/completed/by_completion_date").
					Reply(200).
					File("testdata/complete_task.json")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			output, err := GetTodaysTasks("test-token")
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError {
				if output == "" {
					t.Fatalf("expected output, got empty string")
				}
				if !strings.Contains(output, "Day care pick up") {
					t.Fatalf("output missing task: Day care pick up")
				}
				if !strings.Contains(output, "✅ Take out trash") {
					t.Fatalf("output missing completed task: ✅ Take out trash")
				}
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
				gock.New("https://api.todoist.com").
					Get("/api/v1/tasks/completed/by_completion_date").
					Reply(200).
					File("testdata/complete_task.json")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			output, err := GetCompletedTasks("test-token")
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError {
				if !strings.Contains(output, "✅ Take out trash") {
					t.Fatalf("output missing completed task: ✅ Take out trash")
				}
			}
		})
	}
}
