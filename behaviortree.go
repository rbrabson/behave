package behave

import (
	"strconv"
	"strings"
)

// Status represents the result of a behavior tree node's execution.
type Status int

const (
	Ready   Status = iota // Node is ready to run
	Running               // Node is running
	Success               // Node completed successfully
	Failure               // Node failed
)

// String returns the string representation of the Status.
func (s Status) String() string {
	switch s {
	case Ready:
		return "Ready"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	default:
		return "Unknown"
	}
}

// Node is the interface for all behavior tree nodes.
type Node interface {
	Tick() Status   // Run the node on each tick
	Reset() Status  // Reset the node to initial state
	Status() Status // Get the current status of the node
	String() string // Get a string representation of the node
}

// BehaviorTree represents a behavior tree with a root node.
type BehaviorTree struct {
	Root   Node
	status Status
}

// New creates a new BehaviorTree with the given root node.
func New(root Node) *BehaviorTree {
	return &BehaviorTree{Root: root, status: Ready}
}

// Tick executes the behavior tree.
func (bt *BehaviorTree) Tick() Status {
	if bt.Root == nil {
		bt.status = Failure
		return Failure
	}
	bt.status = bt.Root.Tick()
	return bt.status
}

// Reset resets the behavior tree to its initial state.
func (bt *BehaviorTree) Reset() *BehaviorTree {
	if bt.Root != nil {
		bt.Root.Reset()
	}
	bt.status = Ready
	return bt
}

// Status returns the current status of the behavior tree.
func (bt *BehaviorTree) Status() Status {
	return bt.status
}

// String returns a string representation of the behavior tree.
func (bt *BehaviorTree) String() string {
	var builder strings.Builder
	builder.WriteString("BehaviorTree (" + bt.Status().String() + ")")
	var printNode func(Node, int)
	printNode = func(n Node, depth int) {
		builder.WriteString("\n")
		for i := 0; i < depth; i++ {
			builder.WriteString("  ")
		}
		switch n := n.(type) {
		case *Action:
			builder.WriteString("Action (" + n.Status().String() + ")")
		case *Condition:
			builder.WriteString("Condition (" + n.Status().String() + ")")
		case *Composite:
			builder.WriteString("Composite (" + n.Status().String() + ")")
			for i := range n.Conditions {
				printNode(n.Conditions[i], depth+1)
			}
			if n.Child != nil {
				printNode(n.Child, depth+1)
			}
		case *Sequence:
			builder.WriteString("Sequence (" + n.Status().String() + ")")
			for _, child := range n.Children {
				printNode(child, depth+1)
			}
		case *Selector:
			builder.WriteString("Selector (" + n.Status().String() + ")")
			for _, child := range n.Children {
				printNode(child, depth+1)
			}
		case *Parallel:
			builder.WriteString("Parallel (" + n.Status().String() + ", MinSuccess: " + strconv.Itoa(n.MinSuccessCount) + ")")
			for _, child := range n.Children {
				printNode(child, depth+1)
			}
		case *Retry:
			builder.WriteString("Retry (" + n.Status().String() + ")")
			if n.Child != nil {
				printNode(n.Child, depth+1)
			}
		case *Repeat:
			builder.WriteString("Repeat (" + n.Status().String() + ")")
			if n.Child != nil {
				printNode(n.Child, depth+1)
			}
		default:
			builder.WriteString("Unknown\n")
		}
	}
	printNode(bt.Root, 0)
	return builder.String()
}

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
