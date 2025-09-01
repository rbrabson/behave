package behave

import "strings"

// Selector is a Node that runs its children in order and succeeds if at least one child succeeds.
type Selector struct {
	Children []Node
	status   Status
}

// Reset resets the Selector node and all its children to their initial state.
func (s *Selector) Reset() Status {
	for _, child := range s.Children {
		child.Reset()
	}
	s.status = Ready
	return s.status
}

// Tick executes the selector and handles all status values.
func (s *Selector) Tick() Status {
	for _, child := range s.Children {
		status := child.Tick()
		switch status {
		case Failure:
			// continue to next child
		case Ready, Running, Success:
			s.status = status
			return s.status
		default:
			s.status = Failure
			return s.status
		}
	}
	s.status = Failure
	return s.status
}

// Status returns the current status of the Selector node.
func (s *Selector) Status() Status {
	return s.status
}

// String returns a string representation of the Selector node.
func (s *Selector) String() string {
	var builder strings.Builder
	builder.WriteString("Selector (" + s.Status().String() + ")")
	for _, child := range s.Children {
		builder.WriteString("\n  ")
		builder.WriteString(child.String())
	}
	return builder.String()
}
