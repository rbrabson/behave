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
