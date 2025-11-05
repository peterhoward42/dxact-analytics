package lib

import validator "github.com/go-playground/validator/v10"

// An EventPayload is the type used for comminuicating DrawExact telemetry events.
type EventPayload struct {
	SchemaVersion int `validate:"max=1,min=0"`
	// The EventULID not only makes an event global unique, but also facilitates
	// sorting events into a time ordered sequence.
	EventULID string `validate:"len=26"`
	// The ProxyUserID is our anonymised DrawExact user-key.
	ProxyUserID string `validate:"uuid4"`
	TimeUTC     string `validate:"datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	// The Visit (count) tells us how many cumulative DrawExact sessions the user had launched at the time of this event.
	Visit int `validate:"max=100000,min=1"`
	// The Event tells us what event occurred. E.g. the user completely training step 3.
	Event string `validate:"max=40,min=4"`
	// Some events require additional parameters to adequately describe them.
	// Each event type conveys those parameters in any way it thinks fit encoded into this string.
	Parameters string `validate:"max=80"`
}

const schemaVersion = 1

// NewEventPayload provides an EventPayload ready to use, and includes validation of the payload with respect
// to the field values you pass in.
func NewEventPayload(
	eventULID string,
	proxyUserId string,
	timeUTC string,
	visit int,
	event string,
	parameters string,
) (*EventPayload, error) {
	payload := EventPayload{
		SchemaVersion: schemaVersion,
		EventULID:     eventULID,
		ProxyUserID:   proxyUserId,
		TimeUTC:       timeUTC,
		Visit:         visit,
		Event:         event,
		Parameters:    parameters,
	}
	return &payload, payload.Validate()
}

// Validate performs validation of this Event.
func (payload *EventPayload) Validate() error {
	return validator.New().Struct(payload)
}
