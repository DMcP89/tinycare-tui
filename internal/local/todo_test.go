package local

import (
	"testing"
)

func TestGetLocalTasks(t *testing.T) {
	// Define test cases as a slice of structs.
	tests := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{
			name:      "ValidFile",
			filePath:  "testdata/sample_todo.txt",
			expectErr: false,
		},
		{
			name:      "NonExistentFile",
			filePath:  "testdata/nonexistent.txt",
			expectErr: true,
		},
		{
			name:      "EmptyFile",
			filePath:  "testdata/empty_todo.txt",
			expectErr: false,
		},
	}

	// Single t.Run loop to execute all tests.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetLocalTasks(tc.filePath)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error for file %q, but got none", tc.filePath)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error for file %q, but got: %v", tc.filePath, err)
				}
			}
		})
	}
}
