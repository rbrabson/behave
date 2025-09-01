package behave

import "strings"

// Action is a leaf node that performs an action.
type Action struct {
	Run    func() Status
	status Status
}

// Tick executes the action's Run function and handles all status values.
func (a *Action) Tick() Status {
	if a.Run == nil {
		a.status = Failure
		return a.status
	}
	status := a.Run()
	switch status {
	case Ready, Running, Success, Failure:
		a.status = status
		return a.status
	default:
		a.status = Failure
		return a.status
	}
}

// Reset resets the Action node to its initial state.
func (a *Action) Reset() Status {
	a.status = Ready
	return a.status
}

// Status returns the current status of the Action node.
func (a *Action) Status() Status {
	return a.status
}

// String returns a string representation of the Action node.
func (a *Action) String() string {
	var builder strings.Builder
	builder.WriteString("Action (" + a.Status().String() + ")")
	return builder.String()
}
