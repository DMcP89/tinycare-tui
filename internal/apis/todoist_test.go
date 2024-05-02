package apis

import (
	"testing"

	"github.com/h2non/gock"
)

func TestGetTodaysTasks(t *testing.T) {
	defer gock.Off()

	reqURL := "https://api.todoist.com/rest/v2/tasks?filter=today"
	gock.New(reqURL).
		Get("/").
		Reply(200).
		File("testdata/today_task.json")
	_, err := GetTodaysTasks("test-token")
	if err != nil {
		t.Fatalf("Error while fetching tasks")
	}
}

func TestGetCompletedTasks(t *testing.T) {
	defer gock.Off()

	reqURL := "https://api.todoist.com/sync/v9/completed/get_all"
	gock.New(reqURL).
		Get("/").
		Reply(200).
		File("testdata/complete_task.json")
	_, err := GetCompletedTasks("test-token")
	if err != nil {
		t.Fatalf("Error while fetching tasks")
	}
}
