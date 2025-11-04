package lib

import (
	"fmt"
	"strings"
	"time"
)

// BuildFullPathForRawEvent provides a bucket path and filename for a single raw event. The path facilitates
// date/time partitioning, and guarantees uniqueness.
//
// Example:
//
// events/y=2025/m=11/d=01/hour=14/<eventUUID>.ndjson.gz
//
// The time encoding provides a well known data-lake / Hive-style partitioning pattern.
func BuildFullPathForRawEvent(timeUTC, eventULID string) (path string, err error) {
	fileName := fmt.Sprintf("%s.ndjson.gz", eventULID)
	year, month, day, hour, err := getTimeCompartmentsAsInts(timeUTC)
	if err != nil {
		return
	}
	dateAndTimeSegment := buildDateAndTimeSegment(year, month, day, hour)
	fullPath := strings.Join([]string{"events", dateAndTimeSegment, fileName}, "/")
	return fullPath, nil
}

// getTimeCompartmentsAsInts parses out the year, month etc from the given timestamp,
// having converted the timestamp to UTC.
func getTimeCompartmentsAsInts(timeUTC string) (year, month, day, hour int, err error) {
	theTime, err := time.Parse(time.RFC3339, timeUTC)
	if err != nil {
		return
	}
	asUTC := theTime.UTC()
	return asUTC.Year(), int(asUTC.Month()), asUTC.Day(), asUTC.Hour(), nil
}

func buildDateAndTimeSegment(year, month, day, hour int) string {
	dateAndTimeSegment := fmt.Sprintf("y=%d/m=%02d/d=%02d/hour=%02d", year,month,day,hour)
	return dateAndTimeSegment
}
