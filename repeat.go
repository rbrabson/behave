package behave

import "strings"

// Repeat represents a decorator node that repeats its child until it fails.
// It returns Running while the child succeeds (and resets it for the next iteration),
// Running while the child is running, and Failure when the child fails.
type Repeat struct {
	Child  Node
	status Status
}

// Tick executes the Repeat node, running its child repeatedly until it fails.
func (rp *Repeat) Tick() Status {
	if rp.Child == nil {
		rp.status = Failure
		return rp.status
	}

	childStatus := rp.Child.Tick()
	switch childStatus {
	case Success:
		// Child succeeded, reset it and continue repeating
		rp.Child.Reset()
		rp.status = Running
		return rp.status
	case Running:
		rp.status = Running
		return rp.status
	case Failure:
		// Child failed, we're done
		rp.status = Failure
		return rp.status
	default:
		rp.status = Failure
		return rp.status
	}
}

// Reset resets the Repeat node and its child to the Ready state.
func (rp *Repeat) Reset() Status {
	rp.status = Ready
	if rp.Child != nil {
		rp.Child.Reset()
	}
	return rp.status
}

// Status returns the current status of the Repeat node.
func (rp *Repeat) Status() Status {
	return rp.status
}

// String returns a string representation of the Repeat node.
func (rp *Repeat) String() string {
	var builder strings.Builder
	builder.WriteString("Repeat (")
	builder.WriteString(rp.Status().String())
	builder.WriteString(")")
	if rp.Child != nil {
		builder.WriteString("\n  ")
		builder.WriteString(rp.Child.String())
	}
	return builder.String()
}
