package utils

import (
	"strings"
	"testing"
)

func TestGetSelfCareAdvice(t *testing.T) {
	testString := GetSelfCareAdvice()

	containsAdvice := false
	containsEmoji := false

	for _, s := range advices {
		if strings.Contains(testString, s) {
			containsAdvice = true
		}
	}

	for _, s := range emoji {
		if strings.Contains(testString, s) {
			containsEmoji = true
		}
	}

	if !containsAdvice && !containsEmoji {
		t.Fatalf("Test Failed")
	}
}
