package visitorimplementations

import (
	"github.com/peterhoward42/dxact-analytics/lib"
	"github.com/sanity-io/litter"
)

// SimpleCounter IS-A visitor.Visitor implementation that does a few trivial counts. (To get this process going)
type SimpleCounter struct {

	// userBehaviours captures event counts and flags on a per user basis.
	// It is keyed on a user's proxyID.
	userBehaviours         map[string]*UserBehaviour
	countRecoverableErrors int
	countFatalErrors       int
}

func NewSimpleCounter() *SimpleCounter {
	return &SimpleCounter{
		userBehaviours: map[string]*UserBehaviour{},
	}
}

func (sc *SimpleCounter) Visit(event lib.EventPayload, path string) (err error) {
	switch event.Event {
	case "evt:recoverable-javascript-error":
		sc.countRecoverableErrors += 1
	case "evt:fatal-javascript-error":
		sc.countFatalErrors += 1
	}

	// Initialise the per-user behaviour record the first time we encounter a user.
	var thisUsersBehaviour *UserBehaviour
	var ok bool
	if thisUsersBehaviour, ok = sc.userBehaviours[event.ProxyUserID]; !ok {
		thisUsersBehaviour = NewUserBehaviour()
		sc.userBehaviours[event.ProxyUserID] = thisUsersBehaviour
	}

	switch {

	case event.Event == "evt:launched":
		thisUsersBehaviour.Launched = true

	case event.Event == "evt:enter-training-cage":
		thisUsersBehaviour.EnteredTrainingCage = true

	case event.Event == "evt:training-completed":
		thisUsersBehaviour.CompletedTrainingCage = true

	case event.Event == "evt:sign-in-started":
		thisUsersBehaviour.SignInStarted = true

	case event.Event == "evt:sign-in-success":
		thisUsersBehaviour.SignInSucceeded = true

	case event.Event == "evt:loaded-example":
		thisUsersBehaviour.LoadedAnExample = true

	case event.Event == "evt:created-new-drawing":
		thisUsersBehaviour.CreatedNewDrawing = true

	case event.Event == "evt:retreived-save-drawing":
		thisUsersBehaviour.ReteivedSavedDrawing = true

	}

	return
}

// Report composes a report to out for this telemetry analyser.
func (sc *SimpleCounter) Report() string {
	report := NewReport()

	report.TotalFatalErrors = sc.countFatalErrors
	report.TotalRecoverableErrors = sc.countRecoverableErrors

	howManyPeopleHave := report.HowManyPeopleHave
	for _, behaviour := range sc.userBehaviours {
		if behaviour.Launched {
			howManyPeopleHave.Launched += 1
		}
		if behaviour.EnteredTrainingCage {
			howManyPeopleHave.EnteredTraining += 1
		}
		if behaviour.CompletedTrainingCage {
			howManyPeopleHave.CompletedTraining += 1
		}
		if behaviour.LoadedAnExample {
			howManyPeopleHave.LoadedAnExample += 1
		}
		if behaviour.SignInStarted {
			howManyPeopleHave.TriedToSignIn += 1
		}
		if behaviour.SignInSucceeded {
			howManyPeopleHave.SucceededSigningIn += 1
		}
		if behaviour.CreatedNewDrawing {
			howManyPeopleHave.CreatedTheirOwnDrawing += 1
		}
		if behaviour.ReteivedSavedDrawing {
			howManyPeopleHave.RetreivedTheirASavedDrawing += 1
		}
	}
	return litter.Sdump(report)
}

// A UserBehaviour value tells you if a given user has done various things at any point in their
// DrawExact usage history. They are an approximate indicator of how far they get in the onboarding
// phase.
type UserBehaviour struct {
	Launched              bool
	EnteredTrainingCage   bool
	CompletedTrainingCage bool
	LoadedAnExample       bool
	SignInStarted         bool
	SignInSucceeded       bool
	CreatedNewDrawing     bool
	ReteivedSavedDrawing  bool
}

func NewUserBehaviour() *UserBehaviour {
	return &UserBehaviour{}
}

// HowManyPeopleHave holds a count of how many users have reached various
// onboarding milestones.
type HowManyPeopleHave struct {
	Launched                    int
	EnteredTraining             int
	CompletedTraining           int
	LoadedAnExample             int
	TriedToSignIn               int
	SucceededSigningIn          int
	CreatedTheirOwnDrawing      int
	RetreivedTheirASavedDrawing int
}

type Report struct {
	HowManyPeopleHave      *HowManyPeopleHave
	TotalRecoverableErrors int
	TotalFatalErrors       int
}

func NewReport() *Report {
	return &Report{
		HowManyPeopleHave: &HowManyPeopleHave{},
	}
}
