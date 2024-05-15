package local

import (
	"bufio"
	"fmt"
	"os"
)

func GetLocalTasks(todoFile string) (string, error) {
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
}
