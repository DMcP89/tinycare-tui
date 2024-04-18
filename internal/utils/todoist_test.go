package utils

import (
	"os"
	"testing"
)

func TestGetLocalTasks(t *testing.T) {
	os.Setenv("TODO_FILE", "sample_todo.txt")
	_, err := GetLocalTasks()
	if err != nil {
		t.Fatalf("Error processing sample_todo.txt")
	}
}

func TestGetTodaysTasks(t *testing.T) {
	t.Fatalf("Test not implemented")
}

func TestGetCompletedTasks(t *testing.T) {
	t.Fatalf("Test not implemented")
}
