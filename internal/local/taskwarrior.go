package local

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type TaskwarriorTask struct {
	Description string `json:"description"`
	Status      string `json:"status"`
	Due         string `json:"due"`
	End         string `json:"end"`
}

func GetTaskwarriorTasks() (string, error) {
	out, err := exec.Command("task", "rc.json.array=off", "export").Output()
	if err != nil {
		return "", fmt.Errorf("unable to run task command, is Taskwarrior installed? %w", err)
	}
	tasks, err := parseTaskwarriorExport(out)
	if err != nil {
		return "", err
	}
	return formatTaskwarriorTasks(tasks, time.Now()), nil
}

func parseTaskwarriorExport(data []byte) ([]TaskwarriorTask, error) {
	var tasks []TaskwarriorTask
	if err := json.Unmarshal(data, &tasks); err == nil {
		return tasks, nil
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var task TaskwarriorTask
		if err := json.Unmarshal([]byte(line), &task); err == nil {
			tasks = append(tasks, task)
		}
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("unable to parse taskwarrior export data")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func taskwarriorDateToDay(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	t, err := time.Parse("20060102T150405Z0700", raw)
	if err != nil {
		return "", err
	}
	return t.In(time.Local).Format("2006-01-02"), nil
}

func formatTaskwarriorTasks(tasks []TaskwarriorTask, now time.Time) string {
	today := now.In(time.Local).Format("2006-01-02")
	var output string
	for _, task := range tasks {
		switch task.Status {
		case "pending":
			dueDay, err := taskwarriorDateToDay(task.Due)
			if err != nil || dueDay == "" {
				continue
			}
			if dueDay < today {
				output += fmt.Sprintf("[red]☐ %s[white]\n", task.Description)
			} else if dueDay == today {
				output += fmt.Sprintf("☐ %s\n", task.Description)
			}
		case "completed":
			endDay, err := taskwarriorDateToDay(task.End)
			if err != nil || endDay == "" {
				continue
			}
			if endDay == today {
				output += fmt.Sprintf("✅ %s\n", task.Description)
			}
		}
	}
	return output
}
