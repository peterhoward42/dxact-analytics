package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Register a name for the entry point function.
func init() {
	functions.HTTP("InjestEvent", injestEvent)
}

// injectEvent is an HTTP handler function that expects to
// receive an HTTP POST request carrying a JSON payload that
// matches the SamplePayload type.
//
// All it does is to echo the .Msg field from the POSTed data
// as the response.
func injestEvent(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var payload SamplePayload
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	reply := fmt.Sprintf("Received msg: %s", payload.Msg)
	fmt.Fprintln(w, reply)
}

type SamplePayload struct {
	Msg string
}
