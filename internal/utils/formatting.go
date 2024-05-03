package utils

import (
	"fmt"
	"time"
)

func HumanizeDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	if hours >= 24 {
		return fmt.Sprintf("[green]%d d(s) ago[white]", hours/24)
	} else {
		return fmt.Sprintf("[green]%d h(s) ago[white]", hours)
	}
}
