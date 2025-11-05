package visitorimplementations

import (
	"github.com/peterhoward42/dxact-analytics/lib"
	"github.com/sanity-io/litter"
)

// SimpleCounter IS-A visitor.Visitor implementation that does a few trivial counts. (To get this process going)
type SimpleCounter struct {
	// Users captures all the unique user proxy ids found in the set of events.
	Users map[string]bool

	// All unique event types encountered.
	EventTypes map[string]bool
}

func NewSimpleCounter() *SimpleCounter {
	return &SimpleCounter{
		Users: map[string]bool{},
		EventTypes: map[string]bool{},
	}
}

func (sc *SimpleCounter) Visit(event lib.EventPayload, path string) (err error) {
	// Accumulate the user proxy ID's encountered.
	sc.Users[event.ProxyUserID] = true
	sc.EventTypes[event.Event] = true
	return
}

func (sc *SimpleCounter) Report() string {
	return litter.Sdump(sc)
}
