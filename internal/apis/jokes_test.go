package apis

import (
	"testing"
)

func TestGetJoke(t *testing.T) {
	_, err := GetJoke()
	if err != nil {
		t.Fatalf("Error fetching joke")
	}
}
