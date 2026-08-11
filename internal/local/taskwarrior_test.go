package local

import (
	"os"
	"strings"
	"testing"
	"time"
)

func fixedTaskwarriorNow() time.Time {
	return time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
}

func sampleTaskwarriorTasks() []TaskwarriorTask {
	return []TaskwarriorTask{
		{Description: "Pay bills", Status: "pending", Due: "20260115T120000Z"},
		{Description: "Water plants", Status: "pending", Due: "20260114T120000Z"},
		{Description: "Book flights", Status: "pending", Due: "20260116T120000Z"},
		{Description: "Clean keyboard", Status: "pending"},
		{Description: "Take out trash", Status: "completed", End: "20260115T120000Z"},
		{Description: "Mow the lawn", Status: "completed", End: "20260114T120000Z"},
		{Description: "Reply to email", Status: "waiting", Due: "20260115T120000Z"},
		{Description: "Retired task", Status: "deleted", End: "20260113T120000Z"},
	}
}

func TestParseTaskwarriorExport(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		wantCount int
		expectErr bool
	}{
		{
			name:      "Valid JSON array",
			data:      mustReadFixture(t, "testdata/taskwarrior_export.json"),
			wantCount: 8,
			expectErr: false,
		},
		{
			name:      "Valid newline-delimited JSON",
			data:      mustReadFixture(t, "testdata/taskwarrior_export.ndjson"),
			wantCount: 8,
			expectErr: false,
		},
		{
			name:      "Malformed data returns error",
			data:      []byte("not json at all"),
			wantCount: 0,
			expectErr: true,
		},
		{
			name:      "Empty data returns error",
			data:      []byte(""),
			wantCount: 0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := parseTaskwarriorExport(tt.data)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tasks) != tt.wantCount {
				t.Fatalf("expected %d tasks, got %d", tt.wantCount, len(tasks))
			}
		})
	}
}

func TestFormatTaskwarriorTasks(t *testing.T) {
	now := fixedTaskwarriorNow()

	tests := []struct {
		name     string
		tasks    []TaskwarriorTask
		contains []string
		absent   []string
		want     string
	}{
		{
			name:     "Task due today is listed",
			tasks:    []TaskwarriorTask{{Description: "Pay bills", Status: "pending", Due: "20260115T120000Z"}},
			contains: []string{"☐ Pay bills"},
		},
		{
			name:     "Overdue task is marked red",
			tasks:    []TaskwarriorTask{{Description: "Water plants", Status: "pending", Due: "20260114T120000Z"}},
			contains: []string{"[red]☐ Water plants[white]"},
		},
		{
			name:   "Future task is excluded",
			tasks:  []TaskwarriorTask{{Description: "Book flights", Status: "pending", Due: "20260116T120000Z"}},
			absent: []string{"Book flights"},
		},
		{
			name:   "Task without due date is excluded",
			tasks:  []TaskwarriorTask{{Description: "Clean keyboard", Status: "pending"}},
			absent: []string{"Clean keyboard"},
		},
		{
			name:     "Task completed today is checked",
			tasks:    []TaskwarriorTask{{Description: "Take out trash", Status: "completed", End: "20260115T120000Z"}},
			contains: []string{"✅ Take out trash"},
		},
		{
			name:   "Task completed yesterday is excluded",
			tasks:  []TaskwarriorTask{{Description: "Mow the lawn", Status: "completed", End: "20260114T120000Z"}},
			absent: []string{"Mow the lawn"},
		},
		{
			name:   "Waiting task is excluded",
			tasks:  []TaskwarriorTask{{Description: "Reply to email", Status: "waiting", Due: "20260115T120000Z"}},
			absent: []string{"Reply to email"},
		},
		{
			name:   "Deleted task is excluded",
			tasks:  []TaskwarriorTask{{Description: "Retired task", Status: "deleted", End: "20260113T120000Z"}},
			absent: []string{"Retired task"},
		},
		{
			name:     "Mixed tasks keep only today's and overdue",
			tasks:    sampleTaskwarriorTasks(),
			contains: []string{"☐ Pay bills", "[red]☐ Water plants[white]", "✅ Take out trash"},
			absent:   []string{"Book flights", "Clean keyboard", "Mow the lawn", "Reply to email", "Retired task"},
		},
		{
			name:  "Empty task list produces empty output",
			tasks: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := formatTaskwarriorTasks(tt.tasks, now)
			if tt.want != "" {
				if output != tt.want {
					t.Fatalf("expected output %q, got %q", tt.want, output)
				}
				return
			}
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Fatalf("output missing %q:\n%s", expected, output)
				}
			}
			for _, excluded := range tt.absent {
				if strings.Contains(output, excluded) {
					t.Fatalf("output should not contain %q:\n%s", excluded, output)
				}
			}
		})
	}
}

func TestTaskwarriorDateToDay(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		wantDay   string
		expectErr bool
	}{
		{
			name:    "Valid UTC timestamp",
			raw:     "20260115T120000Z",
			wantDay: "2026-01-15",
		},
		{
			name:    "Empty string returns empty day",
			raw:     "",
			wantDay: "",
		},
		{
			name:      "Malformed timestamp returns error",
			raw:       "not-a-date",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day, err := taskwarriorDateToDay(tt.raw)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if day != tt.wantDay {
				t.Fatalf("expected day %q, got %q", tt.wantDay, day)
			}
		})
	}
}

func mustReadFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unable to read fixture %s: %v", path, err)
	}
	return data
}
