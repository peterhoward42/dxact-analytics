package lib

// An EventPayload instance is the type used for comminuicating DrawExact telemetry events.
type EventPayload struct {
	EventULID   string
	ProxyUserId string
	TimeUTC     string
	Visit       int
	Event       string
	Parameters  string
}

// NewEventPayload provides an EventPayload ready to use.
func NewEventPayload(
	eventULID string,
	proxyUserId string,
	timeUTC string,
	visit int,
	event string,
	parameters string,
) *EventPayload {
	return &EventPayload{
		EventULID:   eventULID,
		ProxyUserId: proxyUserId,
		TimeUTC:     timeUTC,
		Visit:       visit,
		Event:       event,
		Parameters:  parameters,
	}
}
