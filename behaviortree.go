package behave

import (
	"context"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"
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
//
// Parameters:
//   - root: The root node of the behavior tree. This can be any node that implements the Node interface.
//
// Returns:
//   - A pointer to a new BehaviorTree instance initialized with the provided root node and a status of Ready.
func New(root Node) *BehaviorTree {
	return &BehaviorTree{Root: root, status: Ready}
}

// Tick executes the behavior tree.
//
// Returns:
//   - The current status of the behavior tree after execution.
func (bt *BehaviorTree) Tick() Status {
	if bt.Root == nil {
		bt.status = Failure
		return Failure
	}
	bt.status = bt.Root.Tick()
	return bt.status
}

// Reset resets the behavior tree to its initial state.
//
// Returns:
//   - A pointer to the BehaviorTree instance after resetting, allowing for method chaining.
func (bt *BehaviorTree) Reset() *BehaviorTree {
	if bt.Root != nil {
		bt.Root.Reset()
	}
	bt.status = Ready
	return bt
}

// Status returns the current status of the behavior tree.
//
// Returns:
//   - The current status of the behavior tree, which can be Ready, Running, Success, or Failure.
func (bt *BehaviorTree) Status() Status {
	return bt.status
}

// String returns a string representation of the behavior tree.
//
// Returns:
//   - A string that represents the structure and current status of the behavior tree, including its root node and all child nodes.
func (bt *BehaviorTree) String() string {
	var builder strings.Builder
	builder.WriteString("BehaviorTree (" + bt.Status().String() + ")")
	// builder.WriteString("\n. ")
	str := bt.Root.String()
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		builder.WriteString("\n  ")
		builder.WriteString(line)
	}
	return builder.String()
}

// Action is a leaf node that performs an action.
type Action struct {
	Run    func() Status
	status Status
}

// Tick executes the action's Run function and handles all status values.
//
// Returns:
//   - The current status of the Action node after execution, which can be Ready, Running, Success, or Failure.
//     If the Run function is nil or returns an invalid status, it defaults to Failure.
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
//
// Returns:
//   - The status of the Action node after reset, which will be Ready.
func (a *Action) Reset() Status {
	a.status = Ready
	return a.status
}

// Status returns the current status of the Action node.
//
// Returns:
//   - The current status of the Action node, which can be Ready, Running, Success, or Failure.
func (a *Action) Status() Status {
	return a.status
}

// String returns a string representation of the Action node.
//
// Returns:
//   - A string that represents the Action node, including its current status. The format is "Action (Status)".
func (a *Action) String() string {
	var builder strings.Builder
	builder.WriteString("Action (" + a.Status().String() + ")")
	return builder.String()
}

// Condition is a leaf node that checks a condition.
type Condition struct {
	Check func() bool
}

// Tick executes the condition's Check function.
//
// Returns:
//   - Success if the Check function returns true, Failure if it returns false or is nil.
func (c *Condition) Tick() Status {
	return c.Status()
}

// Reset resets the Condition node to its initial state.
//
// Returns:
//   - The status of the Condition node after reset, which will be Ready. However, since Condition nodes are stateless,
//     this method simply returns Ready without changing any internal state.
func (c *Condition) Reset() Status {
	return Ready
}

// Status returns the current status of the Condition node.
//
// Returns:
//   - Success if the Check function returns true, Failure if it returns false or is nil. Since Condition nodes are stateless,
//     this method directly evaluates the Check function each time it's called.
func (c *Condition) Status() Status {
	if c.Check == nil {
		return Failure
	}
	if c.Check() {
		return Success
	}
	return Failure
}

// String returns a string representation of the Condition node.
//
// Returns:
//   - A string that represents the Condition node, including its current status. The format is "Condition (Status)".
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
//
// Returns:
//   - The current status of the Composite node after execution, which can be Ready, Running, Success, or Failure.
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
	for _, condition := range c.Conditions {
		conditionStatus := condition.Tick()
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
//
// Returns:
//   - The status of the Composite node after reset, which will be Ready. This method resets all conditions
//     and the child node (if they exist) to their initial state.
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
//
// Returns:
//   - The current status of the Composite node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the conditions and child node.
func (c *Composite) Status() Status {
	return c.status
}

// String returns a string representation of the Composite node.
//
// Returns:
//   - A string that represents the Composite node, including its current status, all conditions, and the child node (if it exists)
//     The format is: Composite (Status)
func (c *Composite) String() string {
	var builder strings.Builder
	builder.WriteString("Composite (" + c.Status().String() + ")")
	for i := range c.Conditions {
		builder.WriteString("\n  Condition[" + strconv.Itoa(i) + "]: " + c.Conditions[i].String())
	}
	if c.Child != nil {
		childStr := c.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  Child: " + lines[0])
		lines = lines[1:]
		for _, line := range lines {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// Selector is a Node that runs its children in order and succeeds if at least one child succeeds.
// The Selector composite type can be seen as an OR operator with their children.
type Selector struct {
	Children []Node
	status   Status
}

// Reset resets the Selector node and all its children to their initial state.
//
// Returns:
//   - The status of the Selector node after reset, which will be Ready. This method resets all child nodes to their initial state.
func (s *Selector) Reset() Status {
	for _, child := range s.Children {
		child.Reset()
	}
	s.status = Ready
	return s.status
}

// Tick executes the selector and handles all status values.
//
// Returns:
//   - The status of the Selector node after execution, which can be Ready, Running, Success, or Failure.
//     The Selector returns Success if at least one child returns Success, Running if at least one child is
//     Running and none have succeeded, and Failure if all children have failed or are not ready.
func (s *Selector) Tick() Status {
	for _, child := range s.Children {
		status := child.Tick()
		switch status {
		case Failure:
			continue
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
//
// Returns:
//   - The current status of the Selector node, which can be Ready, Running, Success, or Failure.
func (s *Selector) Status() Status {
	return s.status
}

// String returns a string representation of the Selector node.
func (s *Selector) String() string {
	var builder strings.Builder
	builder.WriteString("Selector (" + s.Status().String() + ")")
	for _, child := range s.Children {
		str := child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
	}
	return builder.String()
}

// Sequence is a Node that runs its children in order and succeeds if all children succeed.
// The Sequence composite type can be seen as an AND operator with their children.
type Sequence struct {
	Children []Node
	status   Status
}

// Reset resets the Sequence node and all its children to their initial state.
//
// Returns:
//   - The status of the Sequence node after reset, which will be Ready. This method resets all
//     child nodes to their initial state.
func (s *Sequence) Reset() Status {
	for _, child := range s.Children {
		child.Reset()
	}
	s.status = Ready
	return s.status
}

// Tick runs the sequence and handles all status values.
//
// Returns:
//   - The status of the Sequence node after execution, which can be Ready, Running, Success, or Failure.
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
//
// Returns:
//   - The current status of the Sequence node, which can be Ready, Running, Success, or Failure.
func (s *Sequence) Status() Status {
	return s.status
}

// String returns a string representation of the Sequence node.
//
// Returns:
//   - A string that represents the Sequence node, including its current status and all child nodes.
//     The format is: Sequence (Status)
func (s *Sequence) String() string {
	var builder strings.Builder
	builder.WriteString("Sequence (" + s.Status().String() + ")")
	for _, child := range s.Children {
		str := child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
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
//
// Returns:
//   - The status of the Parallel node after reset, which will be Ready. This method resets all child nodes
//     to their initial state.
func (p *Parallel) Reset() Status {
	for _, child := range p.Children {
		child.Reset()
	}
	p.status = Ready
	return p.status
}

// Tick runs all children in parallel and evaluates based on MinSuccessCount.
//
// Returns:
//   - The status of the Parallel node after execution, which can be Ready, Running, Success, or Failure.
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
//
// Returns:
//   - The current status of the Parallel node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the number of successful,
//     failed, and running children
func (p *Parallel) Status() Status {
	return p.status
}

// String returns a string representation of the Parallel node.
//
// Returns:
//   - A string that represents the Parallel node, including its current status, MinSuccessCount, and all child nodes.
func (p *Parallel) String() string {
	var builder strings.Builder
	builder.WriteString("Parallel (")
	builder.WriteString(p.Status().String())
	builder.WriteString(", MinSuccess: ")
	builder.WriteString(strconv.Itoa(p.MinSuccessCount))
	builder.WriteString(")")
	for _, child := range p.Children {
		str := child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
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
//
// Returns:
//   - The status of the Retry node after execution, which can be Ready, Running, Success, or Failure.
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
//
// Returns:
//   - The status of the Retry node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (r *Retry) Reset() Status {
	r.status = Ready
	if r.Child != nil {
		r.Child.Reset()
	}
	return r.status
}

// Status returns the current status of the Retry node.
//
// Returns:
//   - The current status of the Retry node, which can be Ready, Running, Success, or Failure.
func (r *Retry) Status() Status {
	return r.status
}

// String returns a string representation of the Retry node.
//
// Returns:
//   - A string that represents the Retry node, including its current status and the child node (if it exists).
func (r *Retry) String() string {
	var builder strings.Builder
	builder.WriteString("Retry (")
	builder.WriteString(r.Status().String())
	builder.WriteString(")")
	if r.Child != nil {
		builder.WriteString("\n  ")
		str := r.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
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
//
// Returns:
//   - The status of the Repeat node after execution, which can be Ready, Running, Success, or Failure.
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
//
// Returns:
//   - The status of the Repeat node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (rp *Repeat) Reset() Status {
	rp.status = Ready
	if rp.Child != nil {
		rp.Child.Reset()
	}
	return rp.status
}

// Status returns the current status of the Repeat node.
//
// Returns:
//   - The current status of the Repeat node, which can be Ready, Running, Success, or Failure.
func (rp *Repeat) Status() Status {
	return rp.status
}

// String returns a string representation of the Repeat node.
//
// Returns:
//   - A string that represents the Repeat node, including its current status and the child node (if it exists).
func (rp *Repeat) String() string {
	var builder strings.Builder
	builder.WriteString("Repeat (")
	builder.WriteString(rp.Status().String())
	builder.WriteString(")")
	if rp.Child != nil {
		builder.WriteString("\n  ")
		str := rp.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
	}
	return builder.String()
}

// Invert represents a decorator node that inverts the result of its child.
// It changes Success to Failure and Failure to Success. Running and Ready states pass through unchanged.
type Invert struct {
	Child  Node
	status Status
}

// Tick executes the Invert node, running its child and inverting Success/Failure results.
//
// Returns:
//   - The status of the Invert node after execution, which can be Ready, Running, Success, or Failure.
func (i *Invert) Tick() Status {
	if i.Child == nil {
		i.status = Failure
		return i.status
	}

	childStatus := i.Child.Tick()
	switch childStatus {
	case Success:
		i.status = Failure
		return i.status
	case Failure:
		i.status = Success
		return i.status
	case Running:
		i.status = Running
		return i.status
	case Ready:
		i.status = Ready
		return i.status
	default:
		i.status = Failure
		return i.status
	}
}

// Reset resets the Invert node and its child to the Ready state.
//
// Returns:
//   - The status of the Invert node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (i *Invert) Reset() Status {
	i.status = Ready
	if i.Child != nil {
		i.Child.Reset()
	}
	return i.status
}

// Status returns the current status of the Invert node.
//
// Returns:
//   - The current status of the Invert node, which can be Ready, Running, Success, or Failure.
func (i *Invert) Status() Status {
	return i.status
}

// String returns a string representation of the Invert node.
//
// Returns:
//   - A string that represents the Invert node, including its current status and the child node (if it exists).
func (i *Invert) String() string {
	var builder strings.Builder
	builder.WriteString("Invert (")
	builder.WriteString(i.Status().String())
	builder.WriteString(")")
	if i.Child != nil {
		childStr := i.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  " + lines[0])
		for _, line := range lines[1:] {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// AlwaysSuccess represents a decorator node that always returns Success even if the child fails.
// This is useful for ensuring certain branches always appear successful to their parent nodes.
type AlwaysSuccess struct {
	Child  Node
	status Status
}

// Tick executes the AlwaysSuccess node, running its child but returning Success even if the child fails.
//
// Returns:
//   - The status of the AlwaysSuccess node after execution, which can be Ready, Running, Success, or Failure.
func (as *AlwaysSuccess) Tick() Status {
	if as.Child == nil {
		as.status = Success
		return as.status
	}

	// Execute the child but ignore its result
	as.status = as.Child.Tick()

	// Return Success even if the child failed
	if as.status == Failure {
		as.status = Success
	}

	return as.status
}

// Reset resets the AlwaysSuccess node and its child to the Ready state.
//
// Returns:
//   - The status of the AlwaysSuccess node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (as *AlwaysSuccess) Reset() Status {
	as.status = Ready
	if as.Child != nil {
		as.Child.Reset()
	}
	return as.status
}

// Status returns the current status of the AlwaysSuccess node.
//
// Returns:
//   - The current status of the AlwaysSuccess node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the child node's status but will
//     always return Success if the child failed.
func (as *AlwaysSuccess) Status() Status {
	return as.status
}

// String returns a string representation of the AlwaysSuccess node.
//
// Returns:
//   - A string that represents the AlwaysSuccess node, including its current status and the child node (if it exists).
func (as *AlwaysSuccess) String() string {
	var builder strings.Builder
	builder.WriteString("AlwaysSuccess (")
	builder.WriteString(as.Status().String())
	builder.WriteString(")")
	if as.Child != nil {
		childStr := as.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  " + lines[0])
		for _, line := range lines[1:] {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// AlwaysFailure represents a decorator node that always returns Failure even if the child succeeds.
// This is useful for testing or ensuring certain branches always appear failed to their parent nodes.
type AlwaysFailure struct {
	Child  Node
	status Status
}

// Tick executes the AlwaysFailure node, running its child but returning Failure even if the child succeeds.
//
// Returns:
//   - The status of the AlwaysFailure node after execution, which can be Ready, Running, Success, or Failure.
func (af *AlwaysFailure) Tick() Status {
	if af.Child == nil {
		af.status = Failure
		return af.status
	}

	// Execute the child but ignore its result
	af.status = af.Child.Tick()

	// Return Failure even if the child succeeded
	if af.status == Success {
		af.status = Failure
	}

	return af.status
}

// Reset resets the AlwaysFailure node and its child to the Ready state.
//
// Returns:
//   - The status of the AlwaysFailure node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (af *AlwaysFailure) Reset() Status {
	af.status = Ready
	if af.Child != nil {
		af.Child.Reset()
	}
	return af.status
}

// Status returns the current status of the AlwaysFailure node.
//
// Returns:
//   - The current status of the AlwaysFailure node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the child node's status but will
//     always return Failure if the child succeeded.
func (af *AlwaysFailure) Status() Status {
	return af.status
}

// String returns a string representation of the AlwaysFailure node.
//
// Returns:
//   - A string that represents the AlwaysFailure node, including its current status and the child node (if it exists).
func (af *AlwaysFailure) String() string {
	var builder strings.Builder
	builder.WriteString("AlwaysFailure (")
	builder.WriteString(af.Status().String())
	builder.WriteString(")")
	if af.Child != nil {
		childStr := af.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  " + lines[0])
		for _, line := range lines[1:] {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// RepeatN represents a decorator node that executes its child a specific number of times.
// It returns Running while the execution count is below MaxCount, then returns the child's last result.
type RepeatN struct {
	Child    Node
	MaxCount int // Maximum number of times to execute the child
	Count    int // Current execution count
	status   Status
}

// Tick executes the RepeatN node, running its child up to MaxCount times.
//
// Returns:
//   - The status of the RepeatN node after execution, which can be Ready, Running, Success, or Failure.
//     The node returns Running while the execution count is below MaxCount, and returns the child's last result once MaxCount is reached.
func (rn *RepeatN) Tick() Status {
	if rn.Child == nil {
		rn.status = Failure
		rn.Count = rn.MaxCount // Set count to MaxCount when there's no child
		return rn.status
	}

	// If we haven't reached the maximum count yet
	if rn.MaxCount <= 0 || rn.Count < rn.MaxCount {
		// Execute the child
		childStatus := rn.Child.Tick()

		// If child is still running, don't increment count yet
		if childStatus == Running {
			rn.status = Running
			return rn.status
		}

		// Child completed (Success or Failure), increment count
		rn.Count++

		// If we've reached the maximum count, return the child's result
		if rn.MaxCount > 0 && rn.Count >= rn.MaxCount {
			rn.status = childStatus
			return rn.status
		}

		// We need to run more times, reset the child for the next execution and return Running
		rn.Child.Reset()
		rn.status = Running
		return rn.status
	}

	// We've already completed all executions, return the stored result
	return rn.status
}

// Reset resets the RepeatN node and its child to the Ready state, and resets the execution count.
//
// Returns:
//   - The status of the RepeatN node after reset, which will be Ready. This method resets the child node
//     to its initial state and resets the execution count to zero.
func (rn *RepeatN) Reset() Status {
	rn.status = Ready
	rn.Count = 0
	if rn.Child != nil {
		rn.Child.Reset()
	}
	return rn.status
}

// Status returns the current status of the RepeatN node.
//
// Returns:
//   - The current status of the RepeatN node, which can be Ready, Running, Success, or Failure.
func (rn *RepeatN) Status() Status {
	return rn.status
}

// String returns a string representation of the RepeatN node.
//
// Returns:
//   - A string that represents the RepeatN node, including its current status, execution count, maximum count,
//     and the child node (if it exists).
func (rn *RepeatN) String() string {
	var builder strings.Builder
	builder.WriteString("RepeatN (")
	builder.WriteString(rn.Status().String())
	builder.WriteString(", ")
	builder.WriteString(strconv.Itoa(rn.Count))
	builder.WriteString("/")
	builder.WriteString(strconv.Itoa(rn.MaxCount))
	builder.WriteString(")")
	if rn.Child != nil {
		childStr := rn.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  " + lines[0])
		for _, line := range lines[1:] {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// Forever represents a decorator node that runs its child forever, ignoring its status.
type Forever struct {
	Child  Node
	status Status
}

// Tick executes the Forever node, always returning Running regardless of the child's status.
//
// Returns:
//   - The status of the Forever node after execution, which will always be Running. This node ignores
//     the child's status and continues running indefinitely.
func (f *Forever) Tick() Status {
	if f.Child != nil {
		f.Child.Tick()
	}
	f.status = Running
	return Running
}

// Reset resets the Forever node and its child to the Ready state.
//
// Returns:
//   - The status of the Forever node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (f *Forever) Reset() Status {
	f.status = Ready
	if f.Child != nil {
		f.Child.Reset()
	}
	return f.status
}

// Status returns the current status of the Forever node.
//
// Returns:
//   - The current status of the Forever node, which will always be Running after Tick is called,
//     and Ready after Reset is called.
func (f *Forever) Status() Status {
	return f.status
}

// String returns a string representation of the Forever node.
//
// Returns:
//   - A string that represents the Forever node, including its current status and the child node (if it exists).
func (f *Forever) String() string {
	var builder strings.Builder
	builder.WriteString("Forever (")
	builder.WriteString(f.status.String())
	builder.WriteString(")")
	if f.Child != nil {
		childStr := f.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  " + lines[0])
		for _, line := range lines[1:] {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// WhileSuccess represents a decorator node that returns Running as long as its child
// is either Running or Success, and returns Failure otherwise.
// This is useful for creating loops that continue while a condition remains true.
type WhileSuccess struct {
	Child  Node
	status Status
}

// Tick executes the WhileSuccess node, running its child and continuing while it succeeds or runs.
//
// Returns:
//   - The status of the WhileSuccess node after execution, which can be Ready, Running, Success, or Failure.
//     The node returns Running while the child is Running or Success, and returns Failure if the child fails or is not ready.
func (ws *WhileSuccess) Tick() Status {
	if ws.Child == nil {
		ws.status = Failure
		return ws.status
	}

	// Execute the child
	childStatus := ws.Child.Tick()

	// Continue running if child is Running or Success
	if childStatus == Running || childStatus == Success {
		// If child succeeded, reset it for the next iteration and return Running
		if childStatus == Success {
			ws.Child.Reset()
		}
		ws.status = Running
		return ws.status
	}

	// Child failed (or returned Ready), so we fail
	ws.status = Failure
	return ws.status
}

// Reset resets the WhileSuccess node and its child to the Ready state.
//
// Returns:
//   - The status of the WhileSuccess node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (ws *WhileSuccess) Reset() Status {
	ws.status = Ready
	if ws.Child != nil {
		ws.Child.Reset()
	}
	return ws.status
}

// Status returns the current status of the WhileSuccess node.
//
// Returns:
//   - The current status of the WhileSuccess node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the child node's status.
//     The node returns Running while the child is Running or Success, and returns Failure if the child fails or is not ready.
func (ws *WhileSuccess) Status() Status {
	return ws.status
}

// String returns a string representation of the WhileSuccess node.
//
// Returns:
//   - A string that represents the WhileSuccess node, including its current status and the child node (if it exists).
func (ws *WhileSuccess) String() string {
	var builder strings.Builder
	builder.WriteString("WhileSuccess (")
	builder.WriteString(ws.status.String())
	builder.WriteString(")")
	if ws.Child != nil {
		childStr := ws.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  " + lines[0])
		for _, line := range lines[1:] {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// WhileFailure represents a decorator node that returns Running as long as its child
// returns Running or Failure, and returns Success when the child succeeds.
// It continues executing its child while it fails or is running.
type WhileFailure struct {
	Child  Node
	status Status
}

// Tick executes the WhileFailure node, running its child and continuing while it fails or runs.
//
// Returns:
//   - The status of the WhileFailure node after execution, which can be Ready, Running, Success, or Failure.
//     The node returns Running while the child is Running or Failure, and returns Success if the child succeeds.
func (wf *WhileFailure) Tick() Status {
	if wf.Child == nil {
		wf.status = Success // No child means we're done (child "succeeded")
		return wf.status
	}

	childStatus := wf.Child.Tick()

	switch childStatus {
	case Running, Failure:
		// Continue running while child is running or failing
		wf.status = Running
		// Reset child if it failed so it can try again
		if childStatus == Failure {
			wf.Child.Reset()
		}
		return wf.status
	case Success:
		// Child succeeded, so we're done
		wf.status = Success
		return wf.status
	default:
		// Ready state - should not happen in normal execution
		wf.status = Running
		return wf.status
	}
}

// Reset resets the WhileFailure node and its child to the Ready state.
//
// Returns:
//   - The status of the WhileFailure node after reset, which will be Ready. This method resets the child node
//     to its initial state.
func (wf *WhileFailure) Reset() Status {
	wf.status = Ready
	if wf.Child != nil {
		wf.Child.Reset()
	}
	return wf.status
}

// Status returns the current status of the WhileFailure node.
//
// Returns:
//   - The current status of the WhileFailure node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the child node's status.
//     The node returns Running while the child is Running or Failure, and returns Success if the child succeeds.
func (wf *WhileFailure) Status() Status {
	return wf.status
}

// String returns a string representation of the WhileFailure node.
//
// Returns:
//   - A string that represents the WhileFailure node, including its current status and the child node (if it exists).
func (wf *WhileFailure) String() string {
	var builder strings.Builder
	builder.WriteString("WhileFailure (")
	builder.WriteString(wf.status.String())
	builder.WriteString(")")
	if wf.Child != nil {
		builder.WriteString("\n  ")
		str := wf.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
	}
	return builder.String()
}

// WithTimeout represents a decorator node that runs its child for a maximum duration.
// If the child returns Success or Failure before the duration expires, it returns that status.
// If the duration expires while the child is still Running, it returns Failure.
type WithTimeout struct {
	Child     Node
	Duration  time.Duration
	startTime time.Time
	status    Status
}

// Tick executes the WithTimeout node, running its child and enforcing the timeout.
//
// Returns:
//   - The status of the WithTimeout node after execution, which can be Ready, Running, Success, or Failure.
//     The node returns Success or Failure if the child returns those statuses before the timeout expires,
//     and returns Failure if the timeout expires while the child is still Running.
func (wt *WithTimeout) Tick() Status {
	if wt.Child == nil {
		wt.status = Failure
		return wt.status
	}

	// If this is the first tick, start the timer
	if wt.startTime.IsZero() {
		wt.startTime = time.Now()
	}

	childStatus := wt.Child.Tick()

	switch childStatus {
	case Success, Failure:
		wt.status = childStatus
		return wt.status
	case Running:
		if time.Since(wt.startTime) >= wt.Duration {
			wt.status = Failure // Time's up, child is still running
			return wt.status
		}
		wt.status = Running
		return wt.status
	default:
		wt.status = Failure
		return wt.status
	}
}

// Reset resets the WithTimeout node and its child to the Ready state.
//
// Returns:
//   - The status of the WithTimeout node after reset, which will be Ready. This method resets the child node
//     to its initial state and resets the timer.
func (wt *WithTimeout) Reset() Status {
	wt.status = Ready
	wt.startTime = time.Time{}
	if wt.Child != nil {
		wt.Child.Reset()
	}
	return wt.status
}

// Status returns the current status of the WithTimeout node.
//
// Returns:
//   - The current status of the WithTimeout node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the child node's status
//     and whether the timeout has expired.
func (wt *WithTimeout) Status() Status {
	return wt.status
}

// String returns a string representation of the WithTimeout node.
//
// Returns:
//   - A string that represents the WithTimeout node, including its current status, timeout duration, and the child node (if it exists).
func (wt *WithTimeout) String() string {
	var builder strings.Builder
	builder.WriteString("WithTimeout (")
	builder.WriteString(wt.status.String())
	builder.WriteString(", Duration: ")
	builder.WriteString(wt.Duration.String())
	builder.WriteString(")")
	if wt.Child != nil {
		childStr := wt.Child.String()
		lines := strings.Split(childStr, "\n")
		builder.WriteString("\n  " + lines[0])
		for _, line := range lines[1:] {
			builder.WriteString("\n  " + line)
		}
	}
	return builder.String()
}

// Log represents a decorator node that executes its child and logs the result.
// It's useful for debugging and monitoring behavior tree execution.
type Log struct {
	Child    Node
	Message  string      // Optional custom message for logging
	LogLevel *slog.Level // Optional custom log level. If nil, uses default levels based on child status
	status   Status
}

// Tick executes the Log node, running its child and logging the result.
//
// Returns:
//   - The status of the Log node after execution, which can be Ready, Running, Success, or Failure.
//     The node logs the result of the child execution with the specified message and log level (or defaults based on child status).
func (l *Log) Tick() Status {
	if l.Child == nil {
		l.status = Failure

		// Log with custom level if specified, otherwise use Warn for no child
		logLevel := slog.LevelWarn
		if l.LogLevel != nil {
			logLevel = *l.LogLevel
		}

		slog.Log(context.Background(), logLevel, "Log node has no child", "status", l.status.String())
		return l.status
	}

	// Execute the child
	childStatus := l.Child.Tick()
	l.status = childStatus

	// Log the result with context
	message := l.Message
	if message == "" {
		message = "Log node executed"
	}

	// Determine log level - use custom level if specified, otherwise use defaults based on status
	var logLevel slog.Level
	if l.LogLevel != nil {
		logLevel = *l.LogLevel
	} else {
		// Default log levels based on child status
		switch childStatus {
		case Success:
			logLevel = slog.LevelInfo
		case Failure:
			logLevel = slog.LevelWarn
		case Running, Ready:
			logLevel = slog.LevelDebug
		}
	}

	// Log with the determined level
	slog.Log(context.Background(), logLevel, message,
		"child_status", childStatus.String(),
		"child_type", l.getChildType(),
	)

	return l.status
}

// getChildType returns a string representation of the child node type for logging.
//
// Returns:
//   - A string representing the type of the child node, or "nil" if there is no child.
func (l *Log) getChildType() string {
	if l.Child == nil {
		return "nil"
	}

	t := reflect.TypeOf(l.Child)
	return t.Elem().Name()
}

// Reset resets the Log node and its child to the Ready state.
//
// Returns:
//   - The status of the Log node after reset, which will be Ready. This method resets the child node
//     to its initial state and logs the reset action.
func (l *Log) Reset() Status {
	l.status = Ready
	if l.Child != nil {
		l.Child.Reset()
	}

	// Log reset with custom level if specified, otherwise use Debug
	logLevel := slog.LevelDebug
	if l.LogLevel != nil {
		logLevel = *l.LogLevel
	}

	slog.Log(context.Background(), logLevel, "Log node reset", "message", l.Message)
	return l.status
}

// Status returns the current status of the Log node.
//
// Returns:
//   - The current status of the Log node, which can be Ready, Running, Success, or Failure.
//     This status reflects the result of the last Tick execution, which depends on the child
//     node's status and is logged accordingly.
func (l *Log) Status() Status {
	return l.status
}

// String returns a string representation of the Log node.
//
// Returns:
//   - A string that represents the Log node, including its current status, message, log level, and the child node (if it exists).
func (l *Log) String() string {
	var builder strings.Builder
	builder.WriteString("Log (")
	builder.WriteString(l.status.String())
	if l.Message != "" {
		builder.WriteString(", \"")
		builder.WriteString(l.Message)
		builder.WriteString("\"")
	}
	if l.LogLevel != nil {
		builder.WriteString(", Level:")
		builder.WriteString(l.LogLevel.String())
	}
	builder.WriteString(")")
	if l.Child != nil {
		builder.WriteString("\n  ")
		str := l.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
	}
	return builder.String()
}
