package visitor

import "github.com/peterhoward42/dxact-analytics/lib"

// A Visitor implementation owns the knowledge about how to analyse a set of telemetry events in some particular way.
// It expects to have its Visit method called repeatedly as the set of Events is traversed.
type Visitor interface {
	Visit(event lib.EventPayload, path string) (err error)
	Report() string
}