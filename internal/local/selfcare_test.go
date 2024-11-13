package local

import (
	"strings"
	"testing"
)

func TestGetSelfCareAdvice(t *testing.T) {
	testString := GetSelfCareAdvice()

	containsAny := func(list []string) bool {
		for _, item := range list {
			if strings.Contains(testString, item) {
				return true
			}
		}
		return false
	}

	tests := []struct {
		name    string
		checkFn func() bool
		message string
	}{
		{
			name: "ContainsAdvice",
			checkFn: func() bool {
				return containsAny(advices)
			},
			message: "Expected test string to contain at least one advice, but found none.",
		},
		{
			name: "ContainsEmoji",
			checkFn: func() bool {
				return containsAny(emoji)
			},
			message: "Expected test string to contain at least one emoji, but found none.",
		},
		{
			name: "ContainsEitherAdviceOrEmoji",
			checkFn: func() bool {
				return containsAny(advices) || containsAny(emoji)
			},
			message: "Expected test string to contain at least one advice or emoji, but found neither.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.checkFn() {
				t.Errorf("%s Test string: %q", tc.message, testString)
			}
		})
	}
}
