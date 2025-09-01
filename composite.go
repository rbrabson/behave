package behave

import (
	"strconv"
	"strings"
)

// Composite is a node that combines multiple conditions with any other node.
// It first checks all conditions, and if they all succeed, runs the child node.
type Composite struct {
	Conditions []Node
	Child      Node
	status     Status
}

// Tick executes the composite by first checking all conditions, then running the child node if all conditions succeed.
func (c *Composite) Tick() Status {
	if len(c.Conditions) == 0 && c.Child == nil {
		c.status = Failure
		return c.status
	}

	// If no conditions, just run the child node
	if len(c.Conditions) == 0 {
		if c.Child != nil {
			c.status = c.Child.Tick()
			return c.status
		}
		c.status = Failure
		return c.status
	}

	// Check all conditions first (like a sequence - all must succeed)
	for i := range c.Conditions {
		conditionStatus := c.Conditions[i].Tick()
		switch conditionStatus {
		case Success:
			// This condition succeeded, continue to next condition
			continue
		case Running:
			// This condition is still running
			c.status = Running
			return c.status
		case Failure, Ready:
			// This condition failed or not ready, composite fails
			c.status = Failure
			return c.status
		default:
			c.status = Failure
			return c.status
		}
	}

	// All conditions succeeded, run the child node
	if c.Child != nil {
		c.status = c.Child.Tick()
		return c.status
	}
	c.status = Success
	return c.status
}

// Reset resets the Composite node and its conditions and child node to their initial state.
func (c *Composite) Reset() Status {
	for i := range c.Conditions {
		c.Conditions[i].Reset()
	}
	if c.Child != nil {
		c.Child.Reset()
	}
	c.status = Ready
	return c.status
}

// Status returns the current status of the Composite node.
func (c *Composite) Status() Status {
	return c.status
}

// String returns a string representation of the Composite node.
func (c *Composite) String() string {
	var builder strings.Builder
	builder.WriteString("Composite (" + c.Status().String() + ")")
	for i := range c.Conditions {
		builder.WriteString("\n  Node[" + strconv.Itoa(i) + "]: " + c.Conditions[i].String())
	}
	if c.Child != nil {
		builder.WriteString("\n  Child: " + c.Child.String())
	}
	return builder.String()
}
