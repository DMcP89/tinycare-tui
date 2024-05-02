package local

import (
	"os"
	"testing"
)

func TestGetLocalTasks(t *testing.T) {
	os.Setenv("TODO_FILE", "testdata/sample_todo.txt")
	_, err := GetLocalTasks()
	if err != nil {
		t.Fatalf("Error processing sample_todo.txt")
	}
}
