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

	validator "github.com/go-playground/validator/v10"
)

// Register a name for the entry point function.
func init() {
	functions.HTTP("InjestEvent", injestEvent)
}

// injestEvent is the entry point for receiving POST requests
// with a JSON payload that matches the SamplePayload type.
func injestEvent(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		setPreFlightOptionHeaders(w)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var payload *lib.EventPayload
	var ok bool

	// First a plain unmarshal of the payload.
	if payload, ok = parseJSON(w, r); !ok {
		return
	}

	// Validate the payload for integrity, but also to maybe spot malicious requests.
	if err := validator.New().Struct(payload); err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}

	// Sythesise the unique and deterministic bucket path and filename for this event.
	path, err := lib.BuildFullPathForRawEvent(payload.TimeUTC, payload.EventULID)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	fmt.Printf("XXXX Built this path: %s\n", path)

	// Construct a GCS client

	fmt.Printf("XXXX constructing gcsClient\n")

	gcsContext, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	gcsClient, err := storage.NewClient(gcsContext)
	if err != nil {
		processError(w, http.StatusInternalServerError, err)
		return
	}
	defer gcsClient.Close()

	// Obtain the bucket handle
	fmt.Printf("XXXX obtaining bucket handle\n")
	const telemetryBucket = "drawexact-telemetry"
	gcsBucket := gcsClient.Bucket(telemetryBucket)

	// We compose a pipeline of writers that terminates with one that can write to a GCS bucket.
	fmt.Printf("XXXX construct a bucket writer\n")
	bucketWriter := gcsBucket.Object(path).NewWriter(gcsContext)
	bucketWriter.ContentType = "application/x-ndjson"
	bucketWriter.ContentEncoding = "gzip"

	// Next upstream write stage is gzip compression.
	fmt.Printf("XXXX construct a gzip writer\n")
	gzw := gzip.NewWriter(bucketWriter)

	// Next upstream write stage is JSON encoder (NDJSON format).
	fmt.Printf("XXXX construct a json encoding writer\n")
	enc := json.NewEncoder(gzw)
	enc.SetEscapeHTML(false)

	// So now if we do the JSON encode - the output will be first gzipped, then
	// then written to the storage bucket.
	fmt.Printf("XXXX fire the json encoder and pipeline\n")

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

	fmt.Printf("XXXX injestEvent complete, writing an OK header\n")

	w.WriteHeader(http.StatusOK)
}

func processError(w http.ResponseWriter, statusCode int, err error) {
	errMsg := err.Error()
	fmt.Printf("XXXX in processError, the msg is: %s\n", errMsg)

	fmt.Fprintf(w, "%s", errMsg)
	w.WriteHeader(statusCode)
}

func setPreFlightOptionHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.WriteHeader(http.StatusNoContent)
}

func parseJSON(w http.ResponseWriter, r *http.Request) (payload *lib.EventPayload, ok bool) {
	payload = &lib.EventPayload{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(payload); err != nil {
		processError(w, http.StatusBadRequest, err)
		return nil, false
	}
	return payload, true
}
