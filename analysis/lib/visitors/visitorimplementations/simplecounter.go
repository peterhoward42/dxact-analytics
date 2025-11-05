package visitorimplementations

import (
	"github.com/peterhoward42/dxact-analytics/lib"
	"github.com/sanity-io/litter"
)

// SimpleCounter IS-A visitor.Visitor implementation that does a few trivial counts. (To get this process going)
type SimpleCounter struct {
	// UserList captures all the unique user proxy ids found in the set of events.
	UserList map[string]bool
}

func NewSimpleCounter() *SimpleCounter {
	return &SimpleCounter{
		UserList: map[string]bool{},
	}
}

func (sc *SimpleCounter) Visit(event lib.EventPayload, path string) (err error) {
	// Accumulate the user proxy ID's encountered.
	sc.UserList[event.ProxyUserID] = true
	return
}

func (sc *SimpleCounter) Report() string {
	return litter.Sdump(sc)
}
