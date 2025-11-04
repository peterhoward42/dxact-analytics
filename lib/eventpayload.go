package lib

import validator "github.com/go-playground/validator/v10"

// An EventPayload instance is the type used for comminuicating DrawExact telemetry events.
type EventPayload struct {
	EventULID   string `validate:"len=26"`
	ProxyUserID string `validate:"uuid4"`
	TimeUTC     string `validate:"datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	Visit       int    `validate:"max=100000,min=1"`
	Event       string `validate:"max=40,min=4"`
	Parameters  string `validate:"max=80"`
}

// NewEventPayload provides an EventPayload ready to use.
func NewEventPayload(
	eventULID string,
	proxyUserId string,
	timeUTC string,
	visit int,
	event string,
	parameters string,
) (*EventPayload, error) {
	payload := EventPayload{
		EventULID:   eventULID,
		ProxyUserID: proxyUserId,
		TimeUTC:     timeUTC,
		Visit:       visit,
		Event:       event,
		Parameters:  parameters,
	}
	return &payload, validator.New().Struct(payload)
}
