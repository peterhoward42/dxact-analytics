package function

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/peterhoward42/dxact-analytics/lib"
)

// Register a name for the entry point function.
func init() {
	functions.HTTP("InjestEvent", injestEvent)
}

// injestEvent is the entry point for receiving POST requests
// with a JSON payload that matches the lib.EventPayload type.
// It writes the event to a Google Cloud Storage (GCS) as a single file. I.e. one file per event is stored.
// It uses a hierarchical path that encodes the year/month/day and with a globally unique name.
func injestEvent(w http.ResponseWriter, r *http.Request) {

	// CORS
	if r.Method == http.MethodOptions {
		setPreFlightOptionHeaders(w)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// First a plain unmarshal of the payload.
	var payload *lib.EventPayload
	var ok bool
	if payload, ok = parseJSON(w, r); !ok {
		return
	}

	// Validate the payload for integrity, but also to maybe spot malicious requests.
	if err := payload.Validate(); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}

	// Sythesise the unique and deterministic bucket path and filename for this event.
	path, err := lib.BuildFullPathForRawEvent(payload.TimeUTC, payload.EventULID)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}

	// Construct a GCS client
	gcsContext, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	gcsClient, err := storage.NewClient(gcsContext)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	defer gcsClient.Close()

	// Obtain the bucket handle
	// XXXX todo - fix: this hard code bucket violates 12-factor design principles.
	const telemetryBucket = "drawexact-telemetry"
	gcsBucket := gcsClient.Bucket(telemetryBucket)

	// We compose a pipeline of io.Writer(s) that ends with a writer that can write to a GCS bucket.
	bucketWriter := gcsBucket.Object(path).NewWriter(gcsContext)
	bucketWriter.ContentType = "application/x-ndjson"
	bucketWriter.ContentEncoding = "gzip"

	// Next upstream write stage is gzip compression.
	gzw := gzip.NewWriter(bucketWriter)

	// Final upstream write stage is JSON encoder (NDJSON format).
	enc := json.NewEncoder(gzw)
	enc.SetEscapeHTML(false)

	// So now if we do the JSON encode - the encoded output will be first gzipped, then
	// written to the storage bucket.

	if err := enc.Encode(payload); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	// Force the gzip writer to flush.
	if err := gzw.Close(); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	// Force the bucket writer to flush.
	if err := bucketWriter.Close(); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// processError is a DRY helper that writes the given error's message to
// the given http.ResponseWriter's body, and writes the given HTTP status error code
// the the reponse header.
func processError(w http.ResponseWriter, statusCode int, err error) {
	errMsg := err.Error()
	fmt.Fprintf(w, "%s", errMsg)
	w.WriteHeader(statusCode)
}

// setPreFlightOptionHeaders knows how to set the various required CORS
// headers on the given http.ResponseWriter.
func setPreFlightOptionHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.WriteHeader(http.StatusNoContent)
}

// parseJSON is a simple wrapper that adds some error handling around a json encoder
// in the context of an HTTP request/response.
func parseJSON(w http.ResponseWriter, r *http.Request) (payload *lib.EventPayload, ok bool) {
	payload = &lib.EventPayload{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(payload); err != nil {
		processError(w, http.StatusBadRequest, err)
		return nil, false
	}
	return payload, true
}
