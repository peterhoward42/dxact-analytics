package lib

// An EventPayload instance is the type used for comminuicating events.
type EventPayload struct {
	EventULID   string
	ProxyUserId string
	TimeUTC     string
	Visit       int
	Event       string
	Parameters  string
}
