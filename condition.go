package behave

import "strings"

// Condition is a leaf node that checks a condition.
type Condition struct {
	Check  func() Status
	status Status
}

// Tick executes the condition's Check function.
func (c *Condition) Tick() Status {
	if c.Check == nil {
		c.status = Failure
		return c.status
	}
	status := c.Check()
	switch status {
	case Ready, Running, Success, Failure:
		c.status = status
		return c.status
	default:
		c.status = Failure
		return c.status
	}
}

// Reset resets the Condition node to its initial state.
func (c *Condition) Reset() Status {
	c.status = Ready
	return c.status
}

// Status returns the current status of the Condition node.
func (c *Condition) Status() Status {
	return c.status
}

// String returns a string representation of the Condition node.
func (c *Condition) String() string {
	var builder strings.Builder
	builder.WriteString("Condition (" + c.Status().String() + ")")
	return builder.String()
}
