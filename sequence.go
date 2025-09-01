package behave

import "strings"

// Sequence is a Node that runs its children in order and succeeds if all children succeed.
type Sequence struct {
	Children []Node
	status   Status
}

// Reset resets the Sequence node and all its children to their initial state.
func (s *Sequence) Reset() Status {
	for _, child := range s.Children {
		child.Reset()
	}
	s.status = Ready
	return s.status
}

// Tick runs the sequence and handles all status values.
func (s *Sequence) Tick() Status {
	for _, child := range s.Children {
		status := child.Tick()
		switch status {
		case Success:
			// continue to next child
		case Ready, Running, Failure:
			s.status = status
			return s.status
		default:
			s.status = Failure
			return s.status
		}
	}
	s.status = Success
	return s.status
}

// Status returns the current status of the Sequence node.
func (s *Sequence) Status() Status {
	return s.status
}

// String returns a string representation of the Sequence node.
func (s *Sequence) String() string {
	var builder strings.Builder
	builder.WriteString("Sequence (" + s.Status().String() + ")")
	for _, child := range s.Children {
		builder.WriteString("\n  ")
		builder.WriteString(child.String())
	}
	return builder.String()
}
