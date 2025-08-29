package behave

import (
	"strings"
	"testing"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{Success, "Success"},
		{Failure, "Failure"},
		{Ready, "Ready"},
		{Running, "Running"},
	}

	for _, test := range tests {
		if got := test.status.String(); got != test.expected {
			t.Errorf("Status.String() = %v, want %v", got, test.expected)
		}
	}
}

func TestNew(t *testing.T) {
	action := &Action{}
	bt := New(action)

	if bt.Root != action {
		t.Errorf("New() root = %v, want %v", bt.Root, action)
	}
	if bt.Status() != Ready {
		t.Errorf("New() status = %v, want %v", bt.Status(), Ready)
	}
}

func TestBehaviorTree_Init(t *testing.T) {

}

func TestBehaviorTree_Tick(t *testing.T) {
	action := &Action{
		Run: func() Status { return Success },
	}
	bt := New(action)

	status := bt.Tick()
	if status != Success {
		t.Errorf("BehaviorTree.Tick() = %v, want %v", status, Success)
	}
	if bt.Status() != Success {
		t.Errorf("BehaviorTree.Status() after Tick() = %v, want %v", bt.Status(), Success)
	}
}

func TestBehaviorTree_Stop(t *testing.T) {

}

func TestBehaviorTree_String(t *testing.T) {
	action := &Action{}
	bt := New(action)

	str := bt.String()
	if !strings.Contains(str, "BehaviorTree") {
		t.Errorf("BehaviorTree.String() should contain 'BehaviorTree', got %v", str)
	}
	if !strings.Contains(str, "Action") {
		t.Errorf("BehaviorTree.String() should contain 'Action', got %v", str)
	}
}

func TestAction_Init(t *testing.T) {

}

func TestAction_Tick(t *testing.T) {
	tests := []struct {
		name     string
		action   *Action
		expected Status
	}{
		{
			name:     "no run func",
			action:   &Action{},
			expected: Failure,
		},
		{
			name: "success",
			action: &Action{
				Run: func() Status { return Success },
			},
			expected: Success,
		},
		{
			name: "failure",
			action: &Action{
				Run: func() Status { return Failure },
			},
			expected: Failure,
		},
		{
			name: "running",
			action: &Action{
				Run: func() Status { return Running },
			},
			expected: Running,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := test.action.Tick()
			if status != test.expected {
				t.Errorf("Action.Tick() = %v, want %v", status, test.expected)
			}
			if test.action.Status() != test.expected {
				t.Errorf("Action.Status() after Tick() = %v, want %v", test.action.Status(), test.expected)
			}
		})
	}
}

func TestCondition_Tick(t *testing.T) {
	tests := []struct {
		name      string
		condition *Condition
		expected  Status
	}{
		{
			name:      "no check func",
			condition: &Condition{},
			expected:  Failure,
		},
		{
			name: "success",
			condition: &Condition{
				Check: func() Status { return Success },
			},
			expected: Success,
		},
		{
			name: "failure",
			condition: &Condition{
				Check: func() Status { return Failure },
			},
			expected: Failure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := test.condition.Tick()
			if status != test.expected {
				t.Errorf("Condition.Tick() = %v, want %v", status, test.expected)
			}
			if test.condition.Status() != test.expected {
				t.Errorf("Condition.Status() after Tick() = %v, want %v", test.condition.Status(), test.expected)
			}
		})
	}
}

func TestSequence_Tick(t *testing.T) {
	tests := []struct {
		name     string
		children []Node
		expected Status
	}{
		{
			name:     "empty sequence",
			children: []Node{},
			expected: Success,
		},
		{
			name: "all success",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Success }},
			},
			expected: Success,
		},
		{
			name: "first fails",
			children: []Node{
				&Action{Run: func() Status { return Failure }},
				&Action{Run: func() Status { return Success }},
			},
			expected: Failure,
		},
		{
			name: "second fails",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Failure }},
			},
			expected: Failure,
		},
		{
			name: "first running",
			children: []Node{
				&Action{Run: func() Status { return Running }},
				&Action{Run: func() Status { return Success }},
			},
			expected: Running,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sequence := &Sequence{Children: test.children}
			status := sequence.Tick()
			if status != test.expected {
				t.Errorf("Sequence.Tick() = %v, want %v", status, test.expected)
			}
			if sequence.Status() != test.expected {
				t.Errorf("Sequence.Status() after Tick() = %v, want %v", sequence.Status(), test.expected)
			}
		})
	}
}

func TestSelector_Tick(t *testing.T) {
	tests := []struct {
		name     string
		children []Node
		expected Status
	}{
		{
			name:     "empty selector",
			children: []Node{},
			expected: Failure,
		},
		{
			name: "first succeeds",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Failure }},
			},
			expected: Success,
		},
		{
			name: "second succeeds",
			children: []Node{
				&Action{Run: func() Status { return Failure }},
				&Action{Run: func() Status { return Success }},
			},
			expected: Success,
		},
		{
			name: "all fail",
			children: []Node{
				&Action{Run: func() Status { return Failure }},
				&Action{Run: func() Status { return Failure }},
			},
			expected: Failure,
		},
		{
			name: "first running",
			children: []Node{
				&Action{Run: func() Status { return Running }},
				&Action{Run: func() Status { return Success }},
			},
			expected: Running,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			selector := &Selector{Children: test.children}
			status := selector.Tick()
			if status != test.expected {
				t.Errorf("Selector.Tick() = %v, want %v", status, test.expected)
			}
			if selector.Status() != test.expected {
				t.Errorf("Selector.Status() after Tick() = %v, want %v", selector.Status(), test.expected)
			}
		})
	}
}

func TestComplexBehaviorTree(t *testing.T) {
	// Create a complex behavior tree: Selector with Sequence and Action
	action1 := &Action{Run: func() Status { return Failure }}
	action2 := &Action{Run: func() Status { return Success }}
	action3 := &Action{Run: func() Status { return Success }}

	sequence := &Sequence{Children: []Node{action2, action3}}
	selector := &Selector{Children: []Node{action1, sequence}}
	bt := New(selector)

	// Tick the tree
	status := bt.Tick()
	if status != Success {
		t.Errorf("Complex BehaviorTree.Tick() = %v, want %v", status, Success)
	}

	// Check the tree's status
	if bt.Status() != Success {
		t.Errorf("BehaviorTree.Status() after Tick() = %v, want %v", bt.Status(), Success)
	}

	// Verify the string representation includes all nodes
	str := bt.String()
	expectedParts := []string{"BehaviorTree", "Selector", "Action", "Sequence"}
	for _, part := range expectedParts {
		if !strings.Contains(str, part) {
			t.Errorf("BehaviorTree.String() should contain '%s', got %v", part, str)
		}
	}
}

func TestComposite_Tick(t *testing.T) {
	tests := []struct {
		name      string
		condition Node
		child     Node
		expected  Status
	}{
		{
			name:      "no condition or child",
			condition: nil,
			child:     nil,
			expected:  Failure,
		},
		{
			name:      "no condition, child succeeds",
			condition: nil,
			child:     &Action{Run: func() Status { return Success }},
			expected:  Success,
		},
		{
			name:      "no condition, child fails",
			condition: nil,
			child:     &Action{Run: func() Status { return Failure }},
			expected:  Failure,
		},
		{
			name:      "condition succeeds, no child",
			condition: &Condition{Check: func() Status { return Success }},
			child:     nil,
			expected:  Success,
		},
		{
			name:      "condition succeeds, child succeeds",
			condition: &Condition{Check: func() Status { return Success }},
			child:     &Action{Run: func() Status { return Success }},
			expected:  Success,
		},
		{
			name:      "condition succeeds, child fails",
			condition: &Condition{Check: func() Status { return Success }},
			child:     &Action{Run: func() Status { return Failure }},
			expected:  Failure,
		},
		{
			name:      "condition fails, child not executed",
			condition: &Condition{Check: func() Status { return Failure }},
			child:     &Action{Run: func() Status { return Success }},
			expected:  Failure,
		},
		{
			name:      "condition running",
			condition: &Condition{Check: func() Status { return Running }},
			child:     &Action{Run: func() Status { return Success }},
			expected:  Running,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			composite := &Composite{Condition: test.condition, Child: test.child}
			status := composite.Tick()
			if status != test.expected {
				t.Errorf("Composite.Tick() = %v, want %v", status, test.expected)
			}
			if composite.Status() != test.expected {
				t.Errorf("Composite.Status() after Tick() = %v, want %v", composite.Status(), test.expected)
			}
		})
	}
}

func TestComposite_Reset(t *testing.T) {
	condition := &Condition{Check: func() Status { return Success }}
	child := &Action{Run: func() Status { return Success }}
	composite := &Composite{Condition: condition, Child: child}

	// Set some status first
	composite.Tick()

	// Reset should set status to Ready
	status := composite.Reset()
	if status != Ready {
		t.Errorf("Composite.Reset() = %v, want %v", status, Ready)
	}
	if composite.Status() != Ready {
		t.Errorf("Composite.Status() after Reset() = %v, want %v", composite.Status(), Ready)
	}
}

func TestParallel_Tick(t *testing.T) {
	tests := []struct {
		name            string
		children        []Node
		minSuccessCount int
		expected        Status
	}{
		{
			name:            "empty parallel",
			children:        []Node{},
			minSuccessCount: 1,
			expected:        Success,
		},
		{
			name: "all succeed, need all",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Success }},
			},
			minSuccessCount: 2,
			expected:        Success,
		},
		{
			name: "all succeed, need one",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Success }},
			},
			minSuccessCount: 1,
			expected:        Success,
		},
		{
			name: "one succeeds, need one",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Failure }},
			},
			minSuccessCount: 1,
			expected:        Success,
		},
		{
			name: "one succeeds, need two",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Failure }},
			},
			minSuccessCount: 2,
			expected:        Failure,
		},
		{
			name: "all fail",
			children: []Node{
				&Action{Run: func() Status { return Failure }},
				&Action{Run: func() Status { return Failure }},
			},
			minSuccessCount: 1,
			expected:        Failure,
		},
		{
			name: "one running, one success, need two",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Running }},
			},
			minSuccessCount: 2,
			expected:        Running,
		},
		{
			name: "one running, one failure, need two",
			children: []Node{
				&Action{Run: func() Status { return Failure }},
				&Action{Run: func() Status { return Running }},
			},
			minSuccessCount: 2,
			expected:        Failure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parallel := &Parallel{Children: test.children, MinSuccessCount: test.minSuccessCount}
			status := parallel.Tick()
			if status != test.expected {
				t.Errorf("Parallel.Tick() = %v, want %v", status, test.expected)
			}
			if parallel.Status() != test.expected {
				t.Errorf("Parallel.Status() after Tick() = %v, want %v", parallel.Status(), test.expected)
			}
		})
	}
}

func TestParallel_MinSuccessCount_Validation(t *testing.T) {
	tests := []struct {
		name            string
		children        []Node
		minSuccessCount int
		expectedMin     int
	}{
		{
			name: "zero min success count, should be 1",
			children: []Node{
				&Action{Run: func() Status { return Success }},
				&Action{Run: func() Status { return Success }},
			},
			minSuccessCount: 0,
			expectedMin:     1,
		},
		{
			name: "min success count greater than children, should be clamped",
			children: []Node{
				&Action{Run: func() Status { return Success }},
			},
			minSuccessCount: 5,
			expectedMin:     1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parallel := &Parallel{Children: test.children, MinSuccessCount: test.minSuccessCount}
			parallel.Tick() // This will validate and adjust MinSuccessCount
			if parallel.MinSuccessCount != test.expectedMin {
				t.Errorf("Parallel.MinSuccessCount after Tick() = %v, want %v", parallel.MinSuccessCount, test.expectedMin)
			}
		})
	}
}

func TestParallel_Reset(t *testing.T) {
	children := []Node{
		&Action{Run: func() Status { return Success }},
		&Action{Run: func() Status { return Failure }},
	}
	parallel := &Parallel{Children: children, MinSuccessCount: 1}

	// Set some status first
	parallel.Tick()

	// Reset should set status to Ready
	status := parallel.Reset()
	if status != Ready {
		t.Errorf("Parallel.Reset() = %v, want %v", status, Ready)
	}
	if parallel.Status() != Ready {
		t.Errorf("Parallel.Status() after Reset() = %v, want %v", parallel.Status(), Ready)
	}
}

func TestAction_Reset(t *testing.T) {
	action := &Action{Run: func() Status { return Success }}

	// Set some status first
	action.Tick()

	// Reset should set status to Ready
	status := action.Reset()
	if status != Ready {
		t.Errorf("Action.Reset() = %v, want %v", status, Ready)
	}
	if action.Status() != Ready {
		t.Errorf("Action.Status() after Reset() = %v, want %v", action.Status(), Ready)
	}
}

func TestCondition_Reset(t *testing.T) {
	condition := &Condition{Check: func() Status { return Success }}

	// Set some status first
	condition.Tick()

	// Reset should set status to Ready
	status := condition.Reset()
	if status != Ready {
		t.Errorf("Condition.Reset() = %v, want %v", status, Ready)
	}
	if condition.Status() != Ready {
		t.Errorf("Condition.Status() after Reset() = %v, want %v", condition.Status(), Ready)
	}
}

func TestSequence_Reset(t *testing.T) {
	children := []Node{
		&Action{Run: func() Status { return Success }},
		&Action{Run: func() Status { return Success }},
	}
	sequence := &Sequence{Children: children}

	// Set some status first
	sequence.Tick()

	// Reset should set status to Ready and reset children
	status := sequence.Reset()
	if status != Ready {
		t.Errorf("Sequence.Reset() = %v, want %v", status, Ready)
	}
	if sequence.Status() != Ready {
		t.Errorf("Sequence.Status() after Reset() = %v, want %v", sequence.Status(), Ready)
	}
}

func TestSelector_Reset(t *testing.T) {
	children := []Node{
		&Action{Run: func() Status { return Failure }},
		&Action{Run: func() Status { return Success }},
	}
	selector := &Selector{Children: children}

	// Set some status first
	selector.Tick()

	// Reset should set status to Ready and reset children
	status := selector.Reset()
	if status != Ready {
		t.Errorf("Selector.Reset() = %v, want %v", status, Ready)
	}
	if selector.Status() != Ready {
		t.Errorf("Selector.Status() after Reset() = %v, want %v", selector.Status(), Ready)
	}
}

func TestBehaviorTree_Reset(t *testing.T) {
	action := &Action{Run: func() Status { return Success }}
	bt := New(action)

	// Set some status first
	bt.Tick()

	// Reset should set status to Ready and reset root
	result := bt.Reset()
	if result != bt {
		t.Errorf("BehaviorTree.Reset() should return self")
	}
	if bt.Status() != Ready {
		t.Errorf("BehaviorTree.Status() after Reset() = %v, want %v", bt.Status(), Ready)
	}
	if action.Status() != Ready {
		t.Errorf("Root node status after BehaviorTree.Reset() = %v, want %v", action.Status(), Ready)
	}
}

func TestComplexBehaviorTreeWithNewNodes(t *testing.T) {
	// Create a complex tree using all node types
	// Parallel node with Composite children

	// Composite 1: Check health > 50 AND heal if needed
	healthCheck := &Condition{Check: func() Status { return Success }} // Health > 50
	healAction := &Action{Run: func() Status { return Success }}       // Heal action
	composite1 := &Composite{Condition: healthCheck, Child: healAction}

	// Composite 2: Check enemy nearby AND attack
	enemyCheck := &Condition{Check: func() Status { return Success }} // Enemy nearby
	attackAction := &Action{Run: func() Status { return Success }}    // Attack action
	composite2 := &Composite{Condition: enemyCheck, Child: attackAction}

	// Parallel node: Execute both composites, need at least 1 to succeed
	parallel := &Parallel{
		Children:        []Node{composite1, composite2},
		MinSuccessCount: 1,
	}

	// Main sequence: Move to position, then execute parallel actions
	moveAction := &Action{Run: func() Status { return Success }}
	sequence := &Sequence{Children: []Node{moveAction, parallel}}

	// Create behavior tree
	bt := New(sequence)

	// Test initial state
	if bt.Status() != Ready {
		t.Errorf("Initial BehaviorTree.Status() = %v, want %v", bt.Status(), Ready)
	}

	// Test one tick
	status := bt.Tick()
	if status != Success {
		t.Errorf("BehaviorTree.Tick() = %v, want %v", status, Success)
	}
	if bt.Status() != Success {
		t.Errorf("BehaviorTree.Status() after Tick() = %v, want %v", bt.Status(), Success)
	}

	// Test reset
	bt.Reset()
	if bt.Status() != Ready {
		t.Errorf("BehaviorTree.Status() after Reset() = %v, want %v", bt.Status(), Ready)
	}

	// Verify string representation includes all new node types
	str := bt.String()
	expectedParts := []string{"BehaviorTree", "Sequence", "Action", "Parallel", "Composite", "Condition"}
	for _, part := range expectedParts {
		if !strings.Contains(str, part) {
			t.Errorf("BehaviorTree.String() should contain '%s', got %v", part, str)
		}
	}
}
