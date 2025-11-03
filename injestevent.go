package function

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/peterhoward42/dxact-analytics/lib"
	"github.com/sanity-io/litter"
)

// Register a name for the entry point function.
func init() {
	functions.HTTP("InjestEvent", injestEvent)
}

// injestEvent is the entry point for receiving POST requests
// with a JSON payload that matches the SamplePayload type.
func injestEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// First decode the incoming JSON into a Payload struct, taking advantage
	// of the validation performed by DisallowUnknownFields()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var payload lib.EventPayload
	if err := decoder.Decode(&payload); err != nil {
		processError(w, http.StatusBadRequest, err)
		return
	}

	// Perform some validation with two aims:
	// 1) recognise spurious requests from bad actors.
	// 2) ensure the bucket path that will be generated based on the payload is sensible.

	// XXXX todo

	// Sythesise the unique bucket path and filename for this event.
	path, err := lib.BuildFullPathForRawEvent(payload.TimeUTC, payload.EventULID)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	_ = path
	gcsContext, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	gcsClient, err := storage.NewClient(gcsContext)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	defer gcsClient.Close()

	// Re-encode the payload gzip compressed NDJSON.
	var outputBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&outputBuffer)
	enc := json.NewEncoder(gzipWriter)
	enc.SetEscapeHTML(false) // makes it more readable
	if err := enc.Encode(payload); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	if err := gzipWriter.Close(); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}

	gzippedBytes := outputBuffer.Bytes()
	fmt.Printf("XXXX gzippedBytes: %s\n", litter.Sdump(gzippedBytes))

	w.WriteHeader(http.StatusOK)
}

func processError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", err.Error())
}
