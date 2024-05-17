package local

import (
	"testing"
)

func TestGetLocalTasks(t *testing.T) {
	_, err := GetLocalTasks("testdata/sample_todo.txt")
	if err != nil {
		t.Fatalf("Error processing sample_todo.txt")
	}
}
