package behave

import "strings"

// Retry represents a decorator node that retries its child until it succeeds,
// ignoring all failures. It returns Success when the child succeeds, Running
// while the child is running, and keeps retrying (returning Running) when the
// child fails.
type Retry struct {
	Child  Node
	status Status
}

// Tick executes the Retry node, running its child until it succeeds.
func (r *Retry) Tick() Status {
	if r.Child == nil {
		r.status = Failure
		return r.status
	}

	childStatus := r.Child.Tick()
	switch childStatus {
	case Success:
		r.status = Success
		return r.status
	case Running:
		r.status = Running
		return r.status
	case Failure:
		// Ignore failure, reset child and keep trying
		r.Child.Reset()
		r.status = Running
		return r.status
	default:
		r.status = Running
		return r.status
	}
}

// Reset resets the Retry node and its child to the Ready state.
func (r *Retry) Reset() Status {
	r.status = Ready
	if r.Child != nil {
		r.Child.Reset()
	}
	return r.status
}

// Status returns the current status of the Retry node.
func (r *Retry) Status() Status {
	return r.status
}

// String returns a string representation of the Retry node.
func (r *Retry) String() string {
	var builder strings.Builder
	builder.WriteString("Retry (")
	builder.WriteString(r.Status().String())
	builder.WriteString(")")
	if r.Child != nil {
		builder.WriteString("\n  ")
		builder.WriteString(r.Child.String())
	}
	return builder.String()
}
