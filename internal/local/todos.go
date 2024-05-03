package local

import (
	"bufio"
	"fmt"
	"os"
)

func GetLocalTasks() (string, error) {
	// Check for the existance of the environment variable TODO_FILE
	// If it exists return its contents as a string
	// If it does not exist return "Please set your TODO_FILE variable"
	if todoFile, ok := os.LookupEnv("TODO_FILE"); ok {
		file, err := os.Open(todoFile)
		if err != nil {
			return "", fmt.Errorf("Unable to open %s : %w", todoFile, err)
		}
		defer file.Close()

		var output string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			output += fmt.Sprintf("☐ %s\n", scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("Unable to read %s : %w", todoFile, err)
		}
		return output, nil
	} else {
		return "", fmt.Errorf("No TODO_FILE environment variable set")
	}
}
