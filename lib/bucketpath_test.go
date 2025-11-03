package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTimeCompartmentsAsInts(t *testing.T) {
	// Happy path.
	timeUTC := "2025-11-03T12:24:49Z"
	year, month, day, hour, err := getTimeCompartmentsAsInts(timeUTC)
	require.NoError(t, err)
	require.Equal(t, year, 2025)
	require.Equal(t, month, 11)
	require.Equal(t, day, 3)
	require.Equal(t, hour, 12)

	// Error path.
	malformedTime := "2025-malformed-03T12:24:49Z"
	_, _, _, _, err = getTimeCompartmentsAsInts(malformedTime)
	require.Error(t, err)
	errMsg := err.Error()
	require.Contains(t, errMsg, "cannot parse")
}

func TestBuildDateAndTimeSegment(t *testing.T) {
	timeUTC := "2025-11-03T12:24:49Z"
	year, month, day, hour, _ := getTimeCompartmentsAsInts(timeUTC)
	timeSegment := buildDateAndTimeSegment(year, month, day, hour)
	require.Equal(t, "y=2025/m=11/d=03/hour=12", timeSegment)
}

func TestBuildFullPathForRawEvent(t *testing.T) {
	timeUTC := "2025-11-03T12:24:49Z"
	eventULID := "thiseventulid"

	// Happy path
	path, err := BuildFullPathForRawEvent(timeUTC, eventULID)
	require.NoError(t, err)
	require.Equal(t, "events/y=2025/m=11/d=03/hour=12/thiseventulid.ndjson.gz", path)
}
