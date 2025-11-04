package lib

import validator "github.com/go-playground/validator/v10"

// An EventPayload instance is the type used for comminuicating DrawExact telemetry events.
type EventPayload struct {
	SchemaVersion int `validate:"max=1,min=1"`
	// The EventULID not only makes an event global unique, but also facilitates
	// sorting events into a time ordered sequence.
	EventULID string `validate:"len=26"`
	// The ProxyUserID is our anonymised DrawExact user-key.
	ProxyUserID string `validate:"uuid4"`
	TimeUTC     string `validate:"datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	// The Visit (count) tells us if the user was in their first/second/third... DrawExact session, when the event was sent.
	Visit int `validate:"max=100000,min=1"`
	// Event can be thought of as an enumerated type to describe which event has occurred. E.g. completely training step.
	// However the constants are defined in the dxact-wasm repo and this repo is ignorant of them per se.
	Event string `validate:"max=40,min=4"`
	// Some events require additional parameters to adequately describe them.
	// The event conveys these parameters in the Parameters field - in a way that is allowed to be event
	// specific. There is no contract for its format or encoding.
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
	// Note we are doing a formal validation of the payload here, and propagate validation errors back
	// up the call stack by potentially returning an error value.
	return &payload, validator.New().Struct(payload)
}
