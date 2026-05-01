package parser

import (
	"fmt"
	"time"
)

func ParseTimeToTodayUTC(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	result := time.Date(
		today.Year(), today.Month(), today.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0,
		time.UTC,
	)

	return result, nil
}
