package lib

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/require"
)

// TestValidationOnNew exercises the field value validations built in to NewEventPayload()
// with a realistic payload to make sure the validation constraints are set up right.
func TestValidationOnNew(t *testing.T) {
	eventULID := newULID()
	proxyUserID := uuid.NewString()
	visit := 3
	eventType := "telemetry:created-new-drawing"
	parameters := "param1 param2"

	// Happy path
	_, err := NewEventPayload(
		eventULID,
		proxyUserID,
		time.Now().UTC().Format(time.RFC3339),
		visit,
		eventType,
		parameters)

	require.NoError(t, err)

	// Error path
	proxyUserID = "some garbabe"
	_, err = NewEventPayload(
		eventULID,
		proxyUserID,
		time.Now().UTC().Format(time.RFC3339),
		visit ,
		eventType,
		parameters)
		
	errMsg := err.Error()
	require.Contains(t, errMsg, "Error")
	require.Contains(t, errMsg, "Field")
	require.Contains(t, errMsg, "ProxyUserID")
	require.Contains(t, errMsg, "uuid4")
}

func newULID() string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
