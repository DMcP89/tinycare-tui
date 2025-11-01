package main

import (
	"os"
	"testing"
	"time"
)

func TestGetRefreshInterval(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		expected time.Duration
	}{
		{
			name:     "Default value when env not set",
			setEnv:   false,
			expected: 300 * time.Second,
		},
		{
			name:     "Valid positive value",
			envValue: "60",
			setEnv:   true,
			expected: 60 * time.Second,
		},
		{
			name:     "Another valid value",
			envValue: "120",
			setEnv:   true,
			expected: 120 * time.Second,
		},
		{
			name:     "Invalid negative value returns default",
			envValue: "-10",
			setEnv:   true,
			expected: 300 * time.Second,
		},
		{
			name:     "Invalid zero value returns default",
			envValue: "0",
			setEnv:   true,
			expected: 300 * time.Second,
		},
		{
			name:     "Invalid non-numeric value returns default",
			envValue: "invalid",
			setEnv:   true,
			expected: 300 * time.Second,
		},
		{
			name:     "Empty string returns default",
			envValue: "",
			setEnv:   true,
			expected: 300 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variable before each test
			os.Unsetenv("TINYCARE_REFRESH_INTERVAL")

			if tt.setEnv {
				os.Setenv("TINYCARE_REFRESH_INTERVAL", tt.envValue)
				defer os.Unsetenv("TINYCARE_REFRESH_INTERVAL")
			}

			result := GetRefreshInterval()
			if result != tt.expected {
				t.Errorf("GetRefreshInterval() = %v, want %v", result, tt.expected)
			}
		})
	}
}
