package behave

import "strings"

// Status represents the result of a behavior tree node's execution.
type Status int

const (
	Success      Status = iota // Node completed successfully
	Failure                    // Node failed
	Initializing               // Node is initializing
	Ready                      // Node is ready
	Running                    // Node is running
	Stopping                   // Node is stopping
	Stopped                    // Node has stopped
)

// Node is the interface for all behavior tree nodes.
type Node interface {
	Init() Status   // Initialize the node
	Tick() Status   // Run the node on each tick
	Stop() Status   // Stop the node
	Status() Status // Get the current status of the node
	String() string // Get a string representation of the node
}

// BehaviorTree represents a behavior tree with a root node.
type BehaviorTree struct {
	Root   Node
	status Status
}

// Stop stops the behavior tree and its root node.
func (bt *BehaviorTree) Stop() Status {
	if bt.Root != nil {
		bt.status = bt.Root.Stop()
		return bt.status
	}
	bt.status = Stopped
	return bt.status
}

// New creates a new BehaviorTree with the given root node.
func New(root Node) *BehaviorTree {
	return &BehaviorTree{Root: root, status: Ready}
}

// Init initializes the behavior tree and its root node.
func (bt *BehaviorTree) Init() Status {
	if bt.Root != nil {
		bt.status = bt.Root.Init()
		return bt.status
	}
	bt.status = Ready
	return bt.status
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
		default:
			builder.WriteString("Unknown\n")
		}
	}
	printNode(bt.Root, 0)
	return builder.String()
}

// Action is a leaf node that performs an action.
type Action struct {
	InitFunc func() Status
	Run      func() Status
	StopFunc func() Status
	status   Status
}

// Tick executes the action's Run function and handles all status values.
func (a *Action) Tick() Status {
	if a.Run == nil {
		a.status = Failure
		return a.status
	}
	status := a.Run()
	switch status {
	case Initializing, Ready, Running, Stopping, Stopped, Success, Failure:
		a.status = status
		return a.status
	default:
		a.status = Failure
		return a.status
	}
}

// Init for Action calls InitFunc if set, otherwise returns Ready.
func (a *Action) Init() Status {
	if a.InitFunc != nil {
		a.status = a.InitFunc()
		return a.status
	}
	a.status = Ready
	return a.status
}

// Stop for Action calls StopFunc if set, otherwise returns Stopped.
func (a *Action) Stop() Status {
	if a.StopFunc != nil {
		a.status = a.StopFunc()
		return a.status
	}
	a.status = Stopped
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
	InitFunc func() Status
	Check    func() Status
	StopFunc func() Status
	status   Status
}

// Tick executes the condition's Check function.
func (c *Condition) Tick() Status {
	if c.Check == nil {
		c.status = Failure
		return c.status
	}
	status := c.Check()
	switch status {
	case Initializing, Ready, Running, Stopping, Stopped, Success, Failure:
		c.status = status
		return c.status
	default:
		c.status = Failure
		return c.status
	}
}

// Init for Condition calls InitFunc if set, otherwise returns Ready.
func (c *Condition) Init() Status {
	if c.InitFunc != nil {
		c.status = c.InitFunc()
		return c.status
	}
	c.status = Ready
	return c.status
}

// Stop for Condition calls StopFunc if set, otherwise returns Stopped.
func (c *Condition) Stop() Status {
	if c.StopFunc != nil {
		c.status = c.StopFunc()
		return c.status
	}
	c.status = Stopped
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

// Sequence is a Node that runs its children in order until one fails or is running.
type Sequence struct {
	Children []Node
	status   Status
}

// Init initializes all children of the Sequence node.
func (s *Sequence) Init() Status {
	for _, child := range s.Children {
		st := child.Init()
		if st != Ready && st != Success {
			s.status = st
			return st
		}
	}
	s.status = Ready
	return s.status
}

// Stop stops all children of the Sequence node.
func (s *Sequence) Stop() Status {
	for _, child := range s.Children {
		st := child.Stop()
		if st != Stopped && st != Success {
			s.status = st
			return st
		}
	}
	s.status = Stopped
	return s.status
}

// Tick runs the sequence and handles all status values.
func (s *Sequence) Tick() Status {
	for _, child := range s.Children {
		status := child.Tick()
		switch status {
		case Success:
			// continue to next child
		case Ready, Initializing, Running, Stopping, Stopped, Failure:
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

// Selector is a Node that runs its children in order until one succeeds or is running.
type Selector struct {
	Children []Node
	status   Status
}

// Init initializes all children of the Selector node.
func (s *Selector) Init() Status {
	for _, child := range s.Children {
		st := child.Init()
		if st != Ready && st != Success {
			s.status = st
			return st
		}
	}
	s.status = Ready
	return s.status
}

// Stop stops all children of the Selector node.
func (s *Selector) Stop() Status {
	for _, child := range s.Children {
		st := child.Stop()
		if st != Stopped && st != Success {
			s.status = st
			return st
		}
	}
	s.status = Stopped
	return s.status
}

// Tick executes the selector and handles all status values.
func (s *Selector) Tick() Status {
	for _, child := range s.Children {
		status := child.Tick()
		switch status {
		case Failure:
			// continue to next child
		case Ready, Initializing, Running, Stopping, Stopped, Success:
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

// String returns the string representation of the Status.
func (s Status) String() string {
	switch s {
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	case Initializing:
		return "Initializing"
	case Ready:
		return "Ready"
	case Running:
		return "Running"
	case Stopping:
		return "Stopping"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}
