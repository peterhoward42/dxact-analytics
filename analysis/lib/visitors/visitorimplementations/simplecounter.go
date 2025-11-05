package visitorimplementations

import (
	"github.com/peterhoward42/dxact-analytics/lib"
	"github.com/sanity-io/litter"
)

// SimpleCounter IS-A visitor.Visitor implementation that does a few trivial counts. (To get this process going)
type SimpleCounter struct {

	// userBehaviours captures event counts and flags on a per user basis.
	// It is keyed on a user's proxyID.
	userBehaviours map[string]*UserBehaviour
}

func NewSimpleCounter() *SimpleCounter {
	return &SimpleCounter{
		userBehaviours: map[string]*UserBehaviour{},
	}
}

func (sc *SimpleCounter) Visit(event lib.EventPayload, path string) (err error) {
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

func (sc *SimpleCounter) Report() string {
	howManyPeopleHave := HowManyPeopleHave{}
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
	return litter.Sdump(howManyPeopleHave)
}

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
