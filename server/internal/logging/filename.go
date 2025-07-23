package logging

import (
	"fmt"
	"time"
)

// GenerateLogFileName generates a UTC-based log file name in format log-yyyyWwwd-HHMMSS.log
func GenerateLogFileName() string {
	now := time.Now().UTC()

	// Get ISO week (Monday = 1)
	year, week := now.ISOWeek()
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday = 0 in Go, but we want it to be 7
		weekday = 7
	}

	// Format: log-2025W330-143022.log (2025, week 33, Sunday=0->7, dash, 14:30:22)
	return fmt.Sprintf("log-%dW%02d%d-%02d%02d%02d.log",
		year, week, weekday, now.Hour(), now.Minute(), now.Second())
}
