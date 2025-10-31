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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var payload ExpectedPayload
	err := decoder.Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	reply := fmt.Sprintf("%+v", payload)
	fmt.Fprintln(w, reply)
}

type ExpectedPayload struct {
	ProxyUserId string
	TimeUTC     string
	Visit       int
	Event       string
	Parameters  string
}
