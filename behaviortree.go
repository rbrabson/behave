package behave

import (
	"context"
	"log/slog"
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
func (i *Invert) Reset() Status {
	i.status = Ready
	if i.Child != nil {
		i.Child.Reset()
	}
	return i.status
}

// Status returns the current status of the Invert node.
func (i *Invert) Status() Status {
	return i.status
}

// String returns a string representation of the Invert node.
func (i *Invert) String() string {
	var builder strings.Builder
	builder.WriteString("Invert (")
	builder.WriteString(i.Status().String())
	builder.WriteString(")")
	if i.Child != nil {
		builder.WriteString("\n  ")
		str := i.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
	}
	return builder.String()
}

// AlwaysSuccess represents a decorator node that always returns Success regardless of its child's result.
// This is useful for ensuring certain branches always appear successful to their parent nodes.
type AlwaysSuccess struct {
	Child  Node
	status Status
}

// Tick executes the AlwaysSuccess node, running its child but always returning Success.
func (as *AlwaysSuccess) Tick() Status {
	if as.Child == nil {
		as.status = Success
		return as.status
	}

	// Execute the child but ignore its result
	as.Child.Tick()

	// Always return Success regardless of child's result
	as.status = Success
	return as.status
}

// Reset resets the AlwaysSuccess node and its child to the Ready state.
func (as *AlwaysSuccess) Reset() Status {
	as.status = Ready
	if as.Child != nil {
		as.Child.Reset()
	}
	return as.status
}

// Status returns the current status of the AlwaysSuccess node.
func (as *AlwaysSuccess) Status() Status {
	return as.status
}

// String returns a string representation of the AlwaysSuccess node.
func (as *AlwaysSuccess) String() string {
	var builder strings.Builder
	builder.WriteString("AlwaysSuccess (")
	builder.WriteString(as.Status().String())
	builder.WriteString(")")
	if as.Child != nil {
		builder.WriteString("\n  ")
		str := as.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
		}
	}
	return builder.String()
}

// AlwaysFailure represents a decorator node that always returns Failure regardless of its child's result.
// This is useful for testing or ensuring certain branches always appear failed to their parent nodes.
type AlwaysFailure struct {
	Child  Node
	status Status
}

// Tick executes the AlwaysFailure node, running its child but always returning Failure.
func (af *AlwaysFailure) Tick() Status {
	if af.Child == nil {
		af.status = Failure
		return af.status
	}

	// Execute the child but ignore its result
	af.Child.Tick()

	// Always return Failure regardless of child's result
	af.status = Failure
	return af.status
}

// Reset resets the AlwaysFailure node and its child to the Ready state.
func (af *AlwaysFailure) Reset() Status {
	af.status = Ready
	if af.Child != nil {
		af.Child.Reset()
	}
	return af.status
}

// Status returns the current status of the AlwaysFailure node.
func (af *AlwaysFailure) Status() Status {
	return af.status
}

// String returns a string representation of the AlwaysFailure node.
func (af *AlwaysFailure) String() string {
	var builder strings.Builder
	builder.WriteString("AlwaysFailure (")
	builder.WriteString(af.Status().String())
	builder.WriteString(")")
	if af.Child != nil {
		builder.WriteString("\n  ")
		str := af.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
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
func (rn *RepeatN) Tick() Status {
	if rn.Child == nil {
		rn.status = Failure
		rn.Count = rn.MaxCount // Set count to MaxCount when there's no child
		return rn.status
	}

	// Handle edge case: MaxCount is 0
	if rn.MaxCount <= 0 {
		rn.status = rn.Child.Tick()
		// Don't increment count for edge case where MaxCount is 0
		return rn.status
	}

	// If we haven't reached the maximum count yet
	if rn.Count < rn.MaxCount {
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
		if rn.Count >= rn.MaxCount {
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
func (rn *RepeatN) Reset() Status {
	rn.status = Ready
	rn.Count = 0
	if rn.Child != nil {
		rn.Child.Reset()
	}
	return rn.status
}

// Status returns the current status of the RepeatN node.
func (rn *RepeatN) Status() Status {
	return rn.status
}

// String returns a string representation of the RepeatN node.
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
		builder.WriteString("\n  ")
		str := rn.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
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
func (ws *WhileSuccess) Reset() Status {
	ws.status = Ready
	if ws.Child != nil {
		ws.Child.Reset()
	}
	return ws.status
}

// Status returns the current status of the WhileSuccess node.
func (ws *WhileSuccess) Status() Status {
	return ws.status
}

// String returns a string representation of the WhileSuccess node.
func (ws *WhileSuccess) String() string {
	var builder strings.Builder
	builder.WriteString("WhileSuccess (")
	builder.WriteString(ws.status.String())
	builder.WriteString(")")
	if ws.Child != nil {
		builder.WriteString("\n  ")
		str := ws.Child.String()
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			builder.WriteString("\n  ")
			builder.WriteString(line)
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
func (wf *WhileFailure) Reset() Status {
	wf.status = Ready
	if wf.Child != nil {
		wf.Child.Reset()
	}
	return wf.status
}

// Status returns the current status of the WhileFailure node.
func (wf *WhileFailure) Status() Status {
	return wf.status
}

// String returns a string representation of the WhileFailure node.
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

// Log represents a decorator node that executes its child and logs the result.
// It's useful for debugging and monitoring behavior tree execution.
type Log struct {
	Child    Node
	Message  string      // Optional custom message for logging
	LogLevel *slog.Level // Optional custom log level. If nil, uses default levels based on child status
	status   Status
}

// Tick executes the Log node, running its child and logging the result.
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
		"child_type", l.getChildType())

	return l.status
}

// getChildType returns a string representation of the child node type for logging
func (l *Log) getChildType() string {
	if l.Child == nil {
		return "nil"
	}

	// Use type assertion to get the actual type name
	switch l.Child.(type) {
	case *Action:
		return "Action"
	case *Condition:
		return "Condition"
	case *Sequence:
		return "Sequence"
	case *Selector:
		return "Selector"
	case *Parallel:
		return "Parallel"
	case *Composite:
		return "Composite"
	case *Retry:
		return "Retry"
	case *Repeat:
		return "Repeat"
	case *RepeatN:
		return "RepeatN"
	case *Invert:
		return "Invert"
	case *AlwaysSuccess:
		return "AlwaysSuccess"
	case *AlwaysFailure:
		return "AlwaysFailure"
	case *WhileSuccess:
		return "WhileSuccess"
	case *WhileFailure:
		return "WhileFailure"
	case *Log:
		return "Log"
	default:
		return "Unknown"
	}
}

// Reset resets the Log node and its child to the Ready state.
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
func (l *Log) Status() Status {
	return l.status
}

// String returns a string representation of the Log node.
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
