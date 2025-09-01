package behave

import (
	"strconv"
	"strings"
)

// Parallel is a Node that runs all its children in parallel and returns Success
// if at least M children report Success, where M is specified by MinSuccessCount.
type Parallel struct {
	Children        []Node
	MinSuccessCount int
	status          Status
}

// Reset resets the Parallel node and all its children to their initial state.
func (p *Parallel) Reset() Status {
	for _, child := range p.Children {
		child.Reset()
	}
	p.status = Ready
	return p.status
}

// Tick runs all children in parallel and evaluates based on MinSuccessCount.
func (p *Parallel) Tick() Status {
	if len(p.Children) == 0 {
		p.status = Success
		return p.status
	}

	// Validate MinSuccessCount
	if p.MinSuccessCount <= 0 {
		p.MinSuccessCount = 1
	}
	if p.MinSuccessCount > len(p.Children) {
		p.MinSuccessCount = len(p.Children)
	}

	successCount := 0
	failureCount := 0
	runningCount := 0

	// Tick all children
	for _, child := range p.Children {
		status := child.Tick()
		switch status {
		case Success:
			successCount++
		case Failure:
			failureCount++
		case Running:
			runningCount++
		case Ready:
			// Treat Ready as still processing
			runningCount++
		}
	}

	// Check if we have enough successes
	if successCount >= p.MinSuccessCount {
		p.status = Success
		return p.status
	}

	// Check if we can never reach MinSuccessCount (too many failures)
	maxPossibleSuccesses := successCount + runningCount
	if maxPossibleSuccesses < p.MinSuccessCount {
		p.status = Failure
		return p.status
	}

	// Still have a chance to succeed, keep running
	if runningCount > 0 {
		p.status = Running
		return p.status
	}

	// All children are done but we don't have enough successes
	p.status = Failure
	return p.status
}

// Status returns the current status of the Parallel node.
func (p *Parallel) Status() Status {
	return p.status
}

// String returns a string representation of the Parallel node.
func (p *Parallel) String() string {
	var builder strings.Builder
	builder.WriteString("Parallel (")
	builder.WriteString(p.Status().String())
	builder.WriteString(", MinSuccess: ")
	builder.WriteString(strconv.Itoa(p.MinSuccessCount))
	builder.WriteString(")")
	for _, child := range p.Children {
		builder.WriteString("\n  ")
		builder.WriteString(child.String())
	}
	return builder.String()
}
