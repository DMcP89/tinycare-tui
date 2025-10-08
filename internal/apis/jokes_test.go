package apis

import (
	"testing"

	"github.com/h2non/gock"
)

func TestGetJoke(t *testing.T) {
	defer gock.Off()

	tests := []struct {
		name        string
		mockReply   func()
		expectError bool
	}{
		{
			name: "Valid Joke",
			mockReply: func() {
				gock.New(jokeUrl).
					Get("/").
					Reply(200).
					File("testdata/joke.txt")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReply()
			joke, err := GetJoke()
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError && joke == "" {
				t.Fatalf("expected non-empty joke, got empty string")
			}
		})
	}
}
