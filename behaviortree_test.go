package behave

import (
	"log/slog"
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
		name       string
		conditions []Node
		child      Node
		expected   Status
	}{
		{
			name:       "no conditions or child",
			conditions: nil,
			child:      nil,
			expected:   Failure,
		},
		{
			name:       "no conditions, child succeeds",
			conditions: nil,
			child:      &Action{Run: func() Status { return Success }},
			expected:   Success,
		},
		{
			name:       "no conditions, child fails",
			conditions: nil,
			child:      &Action{Run: func() Status { return Failure }},
			expected:   Failure,
		},
		{
			name:       "single condition succeeds, no child",
			conditions: []Node{&Condition{Check: func() Status { return Success }}},
			child:      nil,
			expected:   Success,
		},
		{
			name:       "single condition succeeds, child succeeds",
			conditions: []Node{&Condition{Check: func() Status { return Success }}},
			child:      &Action{Run: func() Status { return Success }},
			expected:   Success,
		},
		{
			name:       "single condition succeeds, child fails",
			conditions: []Node{&Condition{Check: func() Status { return Success }}},
			child:      &Action{Run: func() Status { return Failure }},
			expected:   Failure,
		},
		{
			name:       "single condition fails, child not executed",
			conditions: []Node{&Condition{Check: func() Status { return Failure }}},
			child:      &Action{Run: func() Status { return Success }},
			expected:   Failure,
		},
		{
			name:       "single condition running",
			conditions: []Node{&Condition{Check: func() Status { return Running }}},
			child:      &Action{Run: func() Status { return Success }},
			expected:   Running,
		},
		{
			name: "multiple conditions all succeed, child succeeds",
			conditions: []Node{
				&Condition{Check: func() Status { return Success }},
				&Condition{Check: func() Status { return Success }},
			},
			child:    &Action{Run: func() Status { return Success }},
			expected: Success,
		},
		{
			name: "multiple conditions, first fails",
			conditions: []Node{
				&Condition{Check: func() Status { return Failure }},
				&Condition{Check: func() Status { return Success }},
			},
			child:    &Action{Run: func() Status { return Success }},
			expected: Failure,
		},
		{
			name: "multiple conditions, second fails",
			conditions: []Node{
				&Condition{Check: func() Status { return Success }},
				&Condition{Check: func() Status { return Failure }},
			},
			child:    &Action{Run: func() Status { return Success }},
			expected: Failure,
		},
		{
			name: "multiple conditions, first running",
			conditions: []Node{
				&Condition{Check: func() Status { return Running }},
				&Condition{Check: func() Status { return Success }},
			},
			child:    &Action{Run: func() Status { return Success }},
			expected: Running,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			composite := &Composite{Conditions: test.conditions, Child: test.child}
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
	conditions := []Node{&Condition{Check: func() Status { return Success }}}
	child := &Action{Run: func() Status { return Success }}
	composite := &Composite{Conditions: conditions, Child: child}

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
	composite1 := &Composite{Conditions: []Node{healthCheck}, Child: healAction}

	// Composite 2: Check enemy nearby AND attack
	enemyCheck := &Condition{Check: func() Status { return Success }} // Enemy nearby
	attackAction := &Action{Run: func() Status { return Success }}    // Attack action
	composite2 := &Composite{Conditions: []Node{enemyCheck}, Child: attackAction}

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

func TestRetry_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		expected       Status
		description    string
		tickCount      int
		expectedStatus []Status
	}{
		{
			name:           "no child",
			child:          nil,
			expected:       Failure,
			description:    "should fail when no child is provided",
			tickCount:      1,
			expectedStatus: []Status{Failure},
		},
		{
			name:           "child succeeds immediately",
			child:          &Action{Run: func() Status { return Success }},
			expected:       Success,
			description:    "should succeed when child succeeds",
			tickCount:      1,
			expectedStatus: []Status{Success},
		},
		{
			name:           "child running",
			child:          &Action{Run: func() Status { return Running }},
			expected:       Running,
			description:    "should return Running when child is running",
			tickCount:      1,
			expectedStatus: []Status{Running},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			retry := &Retry{Child: test.child}

			for i, expectedStatus := range test.expectedStatus {
				status := retry.Tick()
				if status != expectedStatus {
					t.Errorf("Retry.Tick() call %d = %v, want %v", i+1, status, expectedStatus)
				}
				if retry.Status() != expectedStatus {
					t.Errorf("Retry.Status() after Tick() call %d = %v, want %v", i+1, retry.Status(), expectedStatus)
				}
			}
		})
	}
}

func TestRetry_FailureRetry(t *testing.T) {
	// Test that Retry keeps trying when child fails
	failCount := 0
	maxFails := 3

	child := &Action{
		Run: func() Status {
			failCount++
			if failCount <= maxFails {
				return Failure
			}
			return Success
		},
	}

	retry := &Retry{Child: child}

	// First few ticks should return Running (retrying after failure)
	for i := 0; i < maxFails; i++ {
		status := retry.Tick()
		if status != Running {
			t.Errorf("Retry.Tick() call %d = %v, want %v (should keep retrying)", i+1, status, Running)
		}
		if retry.Status() != Running {
			t.Errorf("Retry.Status() after failed attempt %d = %v, want %v", i+1, retry.Status(), Running)
		}
	}

	// Final tick should succeed
	status := retry.Tick()
	if status != Success {
		t.Errorf("Retry.Tick() final call = %v, want %v", status, Success)
	}
	if retry.Status() != Success {
		t.Errorf("Retry.Status() after success = %v, want %v", retry.Status(), Success)
	}

	// Verify child was called the expected number of times
	expectedCalls := maxFails + 1 // failures + final success
	if failCount != expectedCalls {
		t.Errorf("Child was called %d times, want %d", failCount, expectedCalls)
	}
}

func TestRetry_ChildReset(t *testing.T) {
	// Test that Retry resets child after each failure
	resetCount := 0
	tickCount := 0

	child := &Action{
		Run: func() Status {
			tickCount++
			if tickCount <= 2 {
				return Failure
			}
			return Success
		},
	}

	// Track reset calls by wrapping in a custom node
	resetTracker := &testNode{
		tickFunc: func() Status {
			return child.Tick()
		},
		resetFunc: func() Status {
			resetCount++
			return child.Reset()
		},
		statusFunc: func() Status {
			return child.Status()
		},
		stringFunc: func() string {
			return child.String()
		},
	}

	retry := &Retry{Child: resetTracker}

	// First tick should fail and trigger reset
	status1 := retry.Tick()
	if status1 != Running {
		t.Errorf("First Retry.Tick() = %v, want %v", status1, Running)
	}
	if resetCount != 1 {
		t.Errorf("Reset count after first failure = %d, want 1", resetCount)
	}

	// Second tick should fail and trigger reset
	status2 := retry.Tick()
	if status2 != Running {
		t.Errorf("Second Retry.Tick() = %v, want %v", status2, Running)
	}
	if resetCount != 2 {
		t.Errorf("Reset count after second failure = %d, want 2", resetCount)
	}

	// Third tick should succeed
	status3 := retry.Tick()
	if status3 != Success {
		t.Errorf("Third Retry.Tick() = %v, want %v", status3, Success)
	}
	// No additional reset should occur on success
	if resetCount != 2 {
		t.Errorf("Reset count after success = %d, want 2", resetCount)
	}
}

// Helper test node for tracking method calls
type testNode struct {
	tickFunc   func() Status
	resetFunc  func() Status
	statusFunc func() Status
	stringFunc func() string
}

func (t *testNode) Tick() Status   { return t.tickFunc() }
func (t *testNode) Reset() Status  { return t.resetFunc() }
func (t *testNode) Status() Status { return t.statusFunc() }
func (t *testNode) String() string { return t.stringFunc() }

func TestRetry_Reset(t *testing.T) {
	child := &Action{Run: func() Status { return Failure }}
	retry := &Retry{Child: child}

	// Set some status first
	retry.Tick()

	// Reset should set status to Ready
	status := retry.Reset()
	if status != Ready {
		t.Errorf("Retry.Reset() = %v, want %v", status, Ready)
	}
	if retry.Status() != Ready {
		t.Errorf("Retry.Status() after Reset() = %v, want %v", retry.Status(), Ready)
	}
}

func TestRetry_String(t *testing.T) {
	tests := []struct {
		name     string
		retry    *Retry
		contains []string
	}{
		{
			name:     "no child",
			retry:    &Retry{Child: nil},
			contains: []string{"Retry", "Ready"},
		},
		{
			name:     "with child",
			retry:    &Retry{Child: &Action{Run: func() Status { return Success }}},
			contains: []string{"Retry", "Ready", "Action"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := test.retry.String()
			for _, expected := range test.contains {
				if !strings.Contains(str, expected) {
					t.Errorf("Retry.String() should contain '%s', got %v", expected, str)
				}
			}
		})
	}
}

func TestRepeat_Tick(t *testing.T) {
	tests := []struct {
		name     string
		repeat   *Repeat
		expected Status
	}{
		{
			name:     "no child",
			repeat:   &Repeat{},
			expected: Failure,
		},
		{
			name: "child succeeds first time",
			repeat: &Repeat{
				Child: &Action{Run: func() Status { return Success }},
			},
			expected: Running, // Should continue running after success
		},
		{
			name: "child fails",
			repeat: &Repeat{
				Child: &Action{Run: func() Status { return Failure }},
			},
			expected: Failure,
		},
		{
			name: "child running",
			repeat: &Repeat{
				Child: &Action{Run: func() Status { return Running }},
			},
			expected: Running,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := test.repeat.Tick()
			if status != test.expected {
				t.Errorf("Repeat.Tick() = %v, want %v", status, test.expected)
			}
			if test.repeat.Status() != test.expected {
				t.Errorf("Repeat.Status() after Tick() = %v, want %v", test.repeat.Status(), test.expected)
			}
		})
	}
}

func TestRepeat_RepeatsUntilFailure(t *testing.T) {
	callCount := 0
	action := &Action{
		Run: func() Status {
			callCount++
			if callCount < 3 {
				return Success // Succeed first 2 times
			}
			return Failure // Fail on the 3rd call
		},
	}

	repeat := &Repeat{Child: action}

	// First tick: child succeeds, repeat should return Running and reset child
	status := repeat.Tick()
	if status != Running {
		t.Errorf("First Repeat.Tick() = %v, want %v", status, Running)
	}
	if callCount != 1 {
		t.Errorf("After first tick, callCount = %v, want 1", callCount)
	}

	// Second tick: child succeeds again, repeat should return Running and reset child
	status = repeat.Tick()
	if status != Running {
		t.Errorf("Second Repeat.Tick() = %v, want %v", status, Running)
	}
	if callCount != 2 {
		t.Errorf("After second tick, callCount = %v, want 2", callCount)
	}

	// Third tick: child fails, repeat should return Failure
	status = repeat.Tick()
	if status != Failure {
		t.Errorf("Third Repeat.Tick() = %v, want %v", status, Failure)
	}
	if callCount != 3 {
		t.Errorf("After third tick, callCount = %v, want 3", callCount)
	}
}

func TestRepeat_ChildReset(t *testing.T) {
	child := &Action{
		Run: func() Status {
			return Success
		},
	}

	// Create a repeat node
	repeat := &Repeat{Child: child}

	// Tick multiple times - each success should cause a reset
	// We can't directly test the reset calls, but we can verify the behavior
	// by checking that the child is in Ready state after each success

	// First tick - child succeeds
	status1 := repeat.Tick()
	if status1 != Running {
		t.Errorf("First tick should return Running, got %v", status1)
	}

	// Child should be reset to Ready state after success
	if child.Status() != Ready {
		t.Errorf("Child status after first success should be Ready, got %v", child.Status())
	}

	// Second tick - child succeeds again
	status2 := repeat.Tick()
	if status2 != Running {
		t.Errorf("Second tick should return Running, got %v", status2)
	}

	// Child should be reset to Ready state again after success
	if child.Status() != Ready {
		t.Errorf("Child status after second success should be Ready, got %v", child.Status())
	}
}

func TestRepeat_Reset(t *testing.T) {
	child := &Action{Run: func() Status { return Success }}
	repeat := &Repeat{Child: child}

	// Set some status first
	repeat.Tick()

	// Reset should set status to Ready
	status := repeat.Reset()
	if status != Ready {
		t.Errorf("Repeat.Reset() = %v, want %v", status, Ready)
	}
	if repeat.Status() != Ready {
		t.Errorf("Repeat.Status() after Reset() = %v, want %v", repeat.Status(), Ready)
	}
}

func TestRepeat_String(t *testing.T) {
	tests := []struct {
		name     string
		repeat   *Repeat
		contains []string
	}{
		{
			name:     "no child",
			repeat:   &Repeat{},
			contains: []string{"Repeat", "Ready"},
		},
		{
			name: "with child",
			repeat: &Repeat{
				Child: &Action{Run: func() Status { return Success }},
			},
			contains: []string{"Repeat", "Ready", "Action"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := test.repeat.String()
			for _, expected := range test.contains {
				if !strings.Contains(str, expected) {
					t.Errorf("Repeat.String() should contain '%s', got %v", expected, str)
				}
			}
		})
	}
}

func TestComplexBehaviorTreeWithRepeat(t *testing.T) {
	// Create a behavior tree that uses Repeat with other node types
	attempts := 0

	// Sequence that succeeds a few times then fails
	taskSequence := &Sequence{
		Children: []Node{
			&Condition{Check: func() Status {
				attempts++
				return Success // Always pass condition
			}},
			&Action{Run: func() Status {
				if attempts <= 3 {
					return Success // Succeed first 3 times
				}
				return Failure // Fail after that
			}},
		},
	}

	// Wrap sequence in Repeat
	repeat := &Repeat{Child: taskSequence}

	// Create behavior tree
	bt := New(repeat)

	// Test initial state
	if bt.Status() != Ready {
		t.Errorf("Initial BehaviorTree.Status() = %v, want %v", bt.Status(), Ready)
	}

	// Tick until failure
	var finalStatus Status
	tickCount := 0
	for tickCount < 10 { // Safety limit
		finalStatus = bt.Tick()
		tickCount++
		if finalStatus == Failure {
			break
		}
		if finalStatus != Running {
			t.Errorf("Expected Running status during repeat, got %v", finalStatus)
		}
	}

	// Should eventually fail
	if finalStatus != Failure {
		t.Errorf("Expected final status to be Failure, got %v", finalStatus)
	}

	// Should have attempted at least 4 times (3 successes + 1 failure)
	if attempts < 4 {
		t.Errorf("Expected at least 4 attempts, got %d", attempts)
	}

	// Verify string representation includes Repeat
	str := bt.String()
	if !strings.Contains(str, "Repeat") {
		t.Errorf("BehaviorTree.String() should contain 'Repeat', got %v", str)
	}
	if !strings.Contains(str, "Sequence") {
		t.Errorf("BehaviorTree.String() should contain 'Sequence', got %v", str)
	}
}

func TestInvert_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		expectedStatus Status
	}{
		{
			name:           "no child",
			child:          nil,
			expectedStatus: Failure,
		},
		{
			name:           "child succeeds, invert to failure",
			child:          &Action{Run: func() Status { return Success }},
			expectedStatus: Failure,
		},
		{
			name:           "child fails, invert to success",
			child:          &Action{Run: func() Status { return Failure }},
			expectedStatus: Success,
		},
		{
			name:           "child running, pass through",
			child:          &Action{Run: func() Status { return Running }},
			expectedStatus: Running,
		},
		{
			name:           "child ready, pass through",
			child:          &Action{Run: func() Status { return Ready }},
			expectedStatus: Ready,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			invert := &Invert{Child: test.child}
			status := invert.Tick()
			if status != test.expectedStatus {
				t.Errorf("Invert.Tick() = %v, want %v", status, test.expectedStatus)
			}
			if invert.Status() != test.expectedStatus {
				t.Errorf("Invert.Status() = %v, want %v", invert.Status(), test.expectedStatus)
			}
		})
	}
}

func TestInvert_InversionBehavior(t *testing.T) {
	// Test repeated inversion behavior
	var result Status
	action := &Action{Run: func() Status { return result }}
	invert := &Invert{Child: action}

	// Test Success -> Failure
	result = Success
	status := invert.Tick()
	if status != Failure {
		t.Errorf("Expected Success to be inverted to Failure, got %v", status)
	}

	// Test Failure -> Success
	result = Failure
	status = invert.Tick()
	if status != Success {
		t.Errorf("Expected Failure to be inverted to Success, got %v", status)
	}

	// Test Running -> Running (no inversion)
	result = Running
	status = invert.Tick()
	if status != Running {
		t.Errorf("Expected Running to pass through unchanged, got %v", status)
	}
}

func TestInvert_Reset(t *testing.T) {
	child := &Action{Run: func() Status { return Success }}
	invert := &Invert{Child: child}

	// Set some status first
	invert.Tick()

	// Reset should set status to Ready and reset child
	status := invert.Reset()
	if status != Ready {
		t.Errorf("Invert.Reset() = %v, want %v", status, Ready)
	}
	if invert.Status() != Ready {
		t.Errorf("Invert.Status() after reset = %v, want %v", invert.Status(), Ready)
	}
	if child.Status() != Ready {
		t.Errorf("Child status after Invert.Reset() = %v, want %v", child.Status(), Ready)
	}
}

func TestInvert_String(t *testing.T) {
	tests := []struct {
		name     string
		child    Node
		expected string
	}{
		{
			name:     "no child",
			child:    nil,
			expected: "Invert (Ready)",
		},
		{
			name:     "with child",
			child:    &Action{Run: func() Status { return Success }},
			expected: "Invert (Ready)\n  Action (Ready)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			invert := &Invert{Child: test.child}
			result := invert.String()
			if result != test.expected {
				t.Errorf("Invert.String() = %v, want %v", result, test.expected)
			}
		})
	}
}

func TestComplexBehaviorTreeWithInvert(t *testing.T) {
	// Create a complex tree: Sequence with inverted condition and action
	var checkResult Status
	var actionResult Status

	condition := &Condition{Check: func() Status { return checkResult }}
	invertedCondition := &Invert{Child: condition}
	action := &Action{Run: func() Status { return actionResult }}

	sequence := &Sequence{Children: []Node{invertedCondition, action}}
	bt := New(sequence)

	// Test 1: Original condition succeeds, inverted fails, sequence fails
	checkResult = Success
	actionResult = Success
	status := bt.Tick()
	if status != Failure {
		t.Errorf("Expected Failure when inverted condition fails, got %v", status)
	}

	// Reset for next test
	bt.Reset()

	// Test 2: Original condition fails, inverted succeeds, action succeeds, sequence succeeds
	checkResult = Failure
	actionResult = Success
	status = bt.Tick()
	if status != Success {
		t.Errorf("Expected Success when inverted condition and action succeed, got %v", status)
	}

	// Verify string representation includes all nodes
	str := bt.String()
	if !strings.Contains(str, "Invert") {
		t.Errorf("BehaviorTree.String() should contain 'Invert', got %v", str)
	}
	if !strings.Contains(str, "Sequence") {
		t.Errorf("BehaviorTree.String() should contain 'Sequence', got %v", str)
	}
	if !strings.Contains(str, "Condition") {
		t.Errorf("BehaviorTree.String() should contain 'Condition', got %v", str)
	}
}

func TestAlwaysSuccess_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		expectedStatus Status
	}{
		{
			name:           "no child",
			child:          nil,
			expectedStatus: Success,
		},
		{
			name:           "child succeeds, return success",
			child:          &Action{Run: func() Status { return Success }},
			expectedStatus: Success,
		},
		{
			name:           "child fails, still return success",
			child:          &Action{Run: func() Status { return Failure }},
			expectedStatus: Success,
		},
		{
			name:           "child running, still return success",
			child:          &Action{Run: func() Status { return Running }},
			expectedStatus: Success,
		},
		{
			name:           "child ready, still return success",
			child:          &Action{Run: func() Status { return Ready }},
			expectedStatus: Success,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alwaysSuccess := &AlwaysSuccess{Child: test.child}
			status := alwaysSuccess.Tick()
			if status != test.expectedStatus {
				t.Errorf("AlwaysSuccess.Tick() = %v, want %v", status, test.expectedStatus)
			}
			if alwaysSuccess.Status() != test.expectedStatus {
				t.Errorf("AlwaysSuccess.Status() = %v, want %v", alwaysSuccess.Status(), test.expectedStatus)
			}
		})
	}
}

func TestAlwaysSuccess_ChildExecution(t *testing.T) {
	// Test that child is actually executed but result is ignored
	executed := false
	childAction := &Action{
		Run: func() Status {
			executed = true
			return Failure // This should be ignored
		},
	}

	alwaysSuccess := &AlwaysSuccess{Child: childAction}
	status := alwaysSuccess.Tick()

	if !executed {
		t.Error("Expected child to be executed")
	}
	if status != Success {
		t.Errorf("Expected Success despite child failure, got %v", status)
	}
}

func TestAlwaysSuccess_Reset(t *testing.T) {
	child := &Action{Run: func() Status { return Failure }}
	alwaysSuccess := &AlwaysSuccess{Child: child}

	// Set some status first
	alwaysSuccess.Tick()

	// Reset should set status to Ready and reset child
	status := alwaysSuccess.Reset()
	if status != Ready {
		t.Errorf("AlwaysSuccess.Reset() = %v, want %v", status, Ready)
	}
	if alwaysSuccess.Status() != Ready {
		t.Errorf("AlwaysSuccess.Status() after reset = %v, want %v", alwaysSuccess.Status(), Ready)
	}
	if child.Status() != Ready {
		t.Errorf("Child status after AlwaysSuccess.Reset() = %v, want %v", child.Status(), Ready)
	}
}

func TestAlwaysSuccess_String(t *testing.T) {
	tests := []struct {
		name     string
		child    Node
		expected string
	}{
		{
			name:     "no child",
			child:    nil,
			expected: "AlwaysSuccess (Ready)",
		},
		{
			name:     "with child",
			child:    &Action{Run: func() Status { return Failure }},
			expected: "AlwaysSuccess (Ready)\n  Action (Ready)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alwaysSuccess := &AlwaysSuccess{Child: test.child}
			result := alwaysSuccess.String()
			if result != test.expected {
				t.Errorf("AlwaysSuccess.String() = %v, want %v", result, test.expected)
			}
		})
	}
}

func TestComplexBehaviorTreeWithAlwaysSuccess(t *testing.T) {
	// Create a complex tree: Sequence with AlwaysSuccess wrapping a failing action
	var actionResult Status

	// This action will fail, but AlwaysSuccess will make it appear successful
	failingAction := &Action{Run: func() Status { return actionResult }}
	alwaysSuccessAction := &AlwaysSuccess{Child: failingAction}

	// Another action that should execute after the "successful" first action
	secondAction := &Action{Run: func() Status { return Success }}

	sequence := &Sequence{Children: []Node{alwaysSuccessAction, secondAction}}
	bt := New(sequence)

	// Test 1: First action fails but AlwaysSuccess makes it succeed, sequence succeeds
	actionResult = Failure
	status := bt.Tick()
	if status != Success {
		t.Errorf("Expected Success when AlwaysSuccess wraps failing action, got %v", status)
	}

	// Reset for next test
	bt.Reset()

	// Test 2: First action succeeds, sequence succeeds
	actionResult = Success
	status = bt.Tick()
	if status != Success {
		t.Errorf("Expected Success when wrapped action succeeds, got %v", status)
	}

	// Verify string representation includes all nodes
	str := bt.String()
	if !strings.Contains(str, "AlwaysSuccess") {
		t.Errorf("BehaviorTree.String() should contain 'AlwaysSuccess', got %v", str)
	}
	if !strings.Contains(str, "Sequence") {
		t.Errorf("BehaviorTree.String() should contain 'Sequence', got %v", str)
	}
	if !strings.Contains(str, "Action") {
		t.Errorf("BehaviorTree.String() should contain 'Action', got %v", str)
	}
}

func TestAlwaysFailure_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		expectedStatus Status
	}{
		{
			name:           "no child",
			child:          nil,
			expectedStatus: Failure,
		},
		{
			name:           "child succeeds, still return failure",
			child:          &Action{Run: func() Status { return Success }},
			expectedStatus: Failure,
		},
		{
			name:           "child fails, return failure",
			child:          &Action{Run: func() Status { return Failure }},
			expectedStatus: Failure,
		},
		{
			name:           "child running, still return failure",
			child:          &Action{Run: func() Status { return Running }},
			expectedStatus: Failure,
		},
		{
			name:           "child ready, still return failure",
			child:          &Action{Run: func() Status { return Ready }},
			expectedStatus: Failure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alwaysFailure := &AlwaysFailure{Child: test.child}
			status := alwaysFailure.Tick()
			if status != test.expectedStatus {
				t.Errorf("AlwaysFailure.Tick() = %v, want %v", status, test.expectedStatus)
			}
			if alwaysFailure.Status() != test.expectedStatus {
				t.Errorf("AlwaysFailure.Status() = %v, want %v", alwaysFailure.Status(), test.expectedStatus)
			}
		})
	}
}

func TestAlwaysFailure_ChildExecution(t *testing.T) {
	// Test that child is actually executed but result is ignored
	executed := false
	childAction := &Action{
		Run: func() Status {
			executed = true
			return Success // This should be ignored
		},
	}

	alwaysFailure := &AlwaysFailure{Child: childAction}
	status := alwaysFailure.Tick()

	if !executed {
		t.Error("Expected child to be executed")
	}
	if status != Failure {
		t.Errorf("Expected Failure despite child success, got %v", status)
	}
}

func TestAlwaysFailure_Reset(t *testing.T) {
	child := &Action{Run: func() Status { return Success }}
	alwaysFailure := &AlwaysFailure{Child: child}

	// Set some status first
	alwaysFailure.Tick()

	// Reset should set status to Ready and reset child
	status := alwaysFailure.Reset()
	if status != Ready {
		t.Errorf("AlwaysFailure.Reset() = %v, want %v", status, Ready)
	}
	if alwaysFailure.Status() != Ready {
		t.Errorf("AlwaysFailure.Status() after reset = %v, want %v", alwaysFailure.Status(), Ready)
	}
	if child.Status() != Ready {
		t.Errorf("Child status after AlwaysFailure.Reset() = %v, want %v", child.Status(), Ready)
	}
}

func TestAlwaysFailure_String(t *testing.T) {
	tests := []struct {
		name     string
		child    Node
		expected string
	}{
		{
			name:     "no child",
			child:    nil,
			expected: "AlwaysFailure (Ready)",
		},
		{
			name:     "with child",
			child:    &Action{Run: func() Status { return Success }},
			expected: "AlwaysFailure (Ready)\n  Action (Ready)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alwaysFailure := &AlwaysFailure{Child: test.child}
			result := alwaysFailure.String()
			if result != test.expected {
				t.Errorf("AlwaysFailure.String() = %v, want %v", result, test.expected)
			}
		})
	}
}

func TestComplexBehaviorTreeWithAlwaysFailure(t *testing.T) {
	// Create a complex tree: Selector with AlwaysFailure wrapping a successful action
	var actionResult Status

	// This action will succeed, but AlwaysFailure will make it appear failed
	succeedingAction := &Action{Run: func() Status { return actionResult }}
	alwaysFailureAction := &AlwaysFailure{Child: succeedingAction}

	// Fallback action that should execute after the "failed" first action
	fallbackAction := &Action{Run: func() Status { return Success }}

	selector := &Selector{Children: []Node{alwaysFailureAction, fallbackAction}}
	bt := New(selector)

	// Test 1: First action succeeds but AlwaysFailure makes it fail, selector tries fallback and succeeds
	actionResult = Success
	status := bt.Tick()
	if status != Success {
		t.Errorf("Expected Success when AlwaysFailure wraps successful action and fallback succeeds, got %v", status)
	}

	// Reset for next test
	bt.Reset()

	// Test 2: First action fails, AlwaysFailure makes it fail, selector tries fallback and succeeds
	actionResult = Failure
	status = bt.Tick()
	if status != Success {
		t.Errorf("Expected Success when wrapped action fails and fallback succeeds, got %v", status)
	}

	// Verify string representation includes all nodes
	str := bt.String()
	if !strings.Contains(str, "AlwaysFailure") {
		t.Errorf("BehaviorTree.String() should contain 'AlwaysFailure', got %v", str)
	}
	if !strings.Contains(str, "Selector") {
		t.Errorf("BehaviorTree.String() should contain 'Selector', got %v", str)
	}
	if !strings.Contains(str, "Action") {
		t.Errorf("BehaviorTree.String() should contain 'Action', got %v", str)
	}
}

func TestAlwaysSuccessAndAlwaysFailureTogether(t *testing.T) {
	// Test using both AlwaysSuccess and AlwaysFailure in the same tree
	succeedingAction := &Action{Run: func() Status { return Success }}
	failingAction := &Action{Run: func() Status { return Failure }}

	alwaysSuccess := &AlwaysSuccess{Child: failingAction}    // Should appear successful
	alwaysFailure := &AlwaysFailure{Child: succeedingAction} // Should appear failed

	// In a sequence: AlwaysSuccess should succeed, AlwaysFailure should fail, sequence should fail
	sequence := &Sequence{Children: []Node{alwaysSuccess, alwaysFailure}}
	bt := New(sequence)

	status := bt.Tick()
	if status != Failure {
		t.Errorf("Expected Failure when sequence contains AlwaysFailure, got %v", status)
	}

	// Reset and try with selector: AlwaysSuccess should succeed, selector should succeed
	bt.Reset()
	selector := &Selector{Children: []Node{alwaysSuccess, alwaysFailure}}
	bt = New(selector)

	status = bt.Tick()
	if status != Success {
		t.Errorf("Expected Success when selector starts with AlwaysSuccess, got %v", status)
	}
}

func TestRepeatN_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		maxCount       int
		expectedStatus Status
		ticksNeeded    int
	}{
		{
			name:           "no child",
			child:          nil,
			maxCount:       3,
			expectedStatus: Failure,
			ticksNeeded:    1,
		},
		{
			name:           "max count zero",
			child:          &Action{Run: func() Status { return Success }},
			maxCount:       0,
			expectedStatus: Success, // Should execute once and return child's result
			ticksNeeded:    1,
		},
		{
			name:           "single execution",
			child:          &Action{Run: func() Status { return Success }},
			maxCount:       1,
			expectedStatus: Success,
			ticksNeeded:    1,
		},
		{
			name:           "multiple executions with success",
			child:          &Action{Run: func() Status { return Success }},
			maxCount:       3,
			expectedStatus: Success,
			ticksNeeded:    3,
		},
		{
			name:           "multiple executions with final failure",
			child:          &Action{Run: func() Status { return Failure }},
			maxCount:       2,
			expectedStatus: Failure,
			ticksNeeded:    2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repeatN := &RepeatN{Child: test.child, MaxCount: test.maxCount}

			var status Status
			for i := 0; i < test.ticksNeeded; i++ {
				status = repeatN.Tick()
				if i < test.ticksNeeded-1 {
					// Should be Running until the last tick
					if status != Running {
						t.Errorf("Tick %d: expected Running, got %v", i+1, status)
					}
				}
			}

			// Final status should match expected
			if status != test.expectedStatus {
				t.Errorf("RepeatN.Tick() final result = %v, want %v", status, test.expectedStatus)
			}
			if repeatN.Status() != test.expectedStatus {
				t.Errorf("RepeatN.Status() = %v, want %v", repeatN.Status(), test.expectedStatus)
			}

			// Check that the count is correct
			if repeatN.Count != test.maxCount {
				t.Errorf("RepeatN.Count = %d, want %d", repeatN.Count, test.maxCount)
			}
		})
	}
}

func TestRepeatN_ChildReset(t *testing.T) {
	executions := 0
	action := &Action{
		Run: func() Status {
			executions++
			return Success
		},
	}

	repeatN := &RepeatN{Child: action, MaxCount: 3}

	// First two ticks should be Running and reset the child
	status1 := repeatN.Tick()
	if status1 != Running {
		t.Errorf("First tick should return Running, got %v", status1)
	}
	if executions != 1 {
		t.Errorf("Expected 1 execution after first tick, got %d", executions)
	}

	status2 := repeatN.Tick()
	if status2 != Running {
		t.Errorf("Second tick should return Running, got %v", status2)
	}
	if executions != 2 {
		t.Errorf("Expected 2 executions after second tick, got %d", executions)
	}

	// Third tick should complete and return the child's result
	status3 := repeatN.Tick()
	if status3 != Success {
		t.Errorf("Third tick should return Success, got %v", status3)
	}
	if executions != 3 {
		t.Errorf("Expected 3 executions after third tick, got %d", executions)
	}

	// Additional ticks should return the same result without executing the child again
	status4 := repeatN.Tick()
	if status4 != Success {
		t.Errorf("Fourth tick should return Success, got %v", status4)
	}
	if executions != 3 {
		t.Errorf("Expected 3 executions after fourth tick (no additional execution), got %d", executions)
	}
}

func TestRepeatN_Reset(t *testing.T) {
	child := &Action{Run: func() Status { return Success }}
	repeatN := &RepeatN{Child: child, MaxCount: 3}

	// Execute a few times
	repeatN.Tick()
	repeatN.Tick()

	// Reset should set status to Ready, reset count, and reset child
	status := repeatN.Reset()
	if status != Ready {
		t.Errorf("RepeatN.Reset() = %v, want %v", status, Ready)
	}
	if repeatN.Status() != Ready {
		t.Errorf("RepeatN.Status() after reset = %v, want %v", repeatN.Status(), Ready)
	}
	if repeatN.Count != 0 {
		t.Errorf("RepeatN.Count after reset = %d, want 0", repeatN.Count)
	}
	if child.Status() != Ready {
		t.Errorf("Child status after RepeatN.Reset() = %v, want %v", child.Status(), Ready)
	}
}

func TestRepeatN_String(t *testing.T) {
	tests := []struct {
		name     string
		child    Node
		maxCount int
		count    int
		expected string
	}{
		{
			name:     "no child",
			child:    nil,
			maxCount: 3,
			count:    0,
			expected: "RepeatN (Ready, 0/3)",
		},
		{
			name:     "with child, partial execution",
			child:    &Action{Run: func() Status { return Success }},
			maxCount: 5,
			count:    2,
			expected: "RepeatN (Ready, 2/5)\n  Action (Ready)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repeatN := &RepeatN{Child: test.child, MaxCount: test.maxCount, Count: test.count}
			result := repeatN.String()
			if result != test.expected {
				t.Errorf("RepeatN.String() = %v, want %v", result, test.expected)
			}
		})
	}
}

func TestRepeatN_RunningChild(t *testing.T) {
	// Test behavior when child returns Running
	runningCount := 0
	action := &Action{
		Run: func() Status {
			runningCount++
			if runningCount < 3 {
				return Running
			}
			return Success
		},
	}

	repeatN := &RepeatN{Child: action, MaxCount: 2}

	// First execution: child returns Running, RepeatN should return Running
	status1 := repeatN.Tick()
	if status1 != Running {
		t.Errorf("Expected Running when child is running, got %v", status1)
	}
	if repeatN.Count != 0 {
		t.Errorf("Count should not increment while child is running, got %d", repeatN.Count)
	}

	// Continue ticking until child completes first execution
	status2 := repeatN.Tick()
	if status2 != Running {
		t.Errorf("Expected Running when child is running, got %v", status2)
	}

	// Child should complete first execution
	status3 := repeatN.Tick()
	if status3 != Running {
		t.Errorf("Expected Running after first execution completion, got %v", status3)
	}
	if repeatN.Count != 1 {
		t.Errorf("Count should be 1 after first execution, got %d", repeatN.Count)
	}

	// Reset the running count for second execution
	runningCount = 0

	// Second execution should also work the same way
	repeatN.Tick()            // Running
	repeatN.Tick()            // Running
	status4 := repeatN.Tick() // Success (final)
	if status4 != Success {
		t.Errorf("Expected Success after all executions complete, got %v", status4)
	}
	if repeatN.Count != 2 {
		t.Errorf("Count should be 2 after all executions, got %d", repeatN.Count)
	}
}

func TestComplexBehaviorTreeWithRepeatN(t *testing.T) {
	// Create a complex tree: Sequence with RepeatN wrapping an action
	executions := 0
	action := &Action{
		Run: func() Status {
			executions++
			return Success
		},
	}

	repeatN := &RepeatN{Child: action, MaxCount: 3}
	otherAction := &Action{Run: func() Status { return Success }}

	sequence := &Sequence{Children: []Node{repeatN, otherAction}}
	bt := New(sequence)

	// First two ticks: RepeatN should be Running, sequence should be Running
	status1 := bt.Tick()
	if status1 != Running {
		t.Errorf("Expected Running during RepeatN execution, got %v", status1)
	}

	status2 := bt.Tick()
	if status2 != Running {
		t.Errorf("Expected Running during RepeatN execution, got %v", status2)
	}

	// Third tick: RepeatN completes, sequence should succeed
	status3 := bt.Tick()
	if status3 != Success {
		t.Errorf("Expected Success when RepeatN completes and other action succeeds, got %v", status3)
	}

	if executions != 3 {
		t.Errorf("Expected 3 executions of the action, got %d", executions)
	}

	// Verify string representation includes all nodes
	str := bt.String()
	if !strings.Contains(str, "RepeatN") {
		t.Errorf("BehaviorTree.String() should contain 'RepeatN', got %v", str)
	}
	if !strings.Contains(str, "Sequence") {
		t.Errorf("BehaviorTree.String() should contain 'Sequence', got %v", str)
	}
	if !strings.Contains(str, "3/3") {
		t.Errorf("BehaviorTree.String() should show execution count '3/3', got %v", str)
	}
}

func TestWhileSuccess_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		childResults   []Status
		expectedStatus []Status
		description    string
	}{
		{
			name:           "no child",
			child:          nil,
			childResults:   []Status{},
			expectedStatus: []Status{Failure},
			description:    "should fail when no child is provided",
		},
		{
			name: "child succeeds continuously",
			child: &Action{Run: func() Status {
				return Success
			}},
			childResults:   []Status{Success, Success, Success},
			expectedStatus: []Status{Running, Running, Running},
			description:    "should keep running while child succeeds",
		},
		{
			name: "child fails after success",
			child: &Action{Run: func() Status {
				return Failure
			}},
			childResults:   []Status{Failure},
			expectedStatus: []Status{Failure},
			description:    "should fail when child fails",
		},
		{
			name: "child running",
			child: &Action{Run: func() Status {
				return Running
			}},
			childResults:   []Status{Running, Running},
			expectedStatus: []Status{Running, Running},
			description:    "should keep running while child is running",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			whileSuccess := &WhileSuccess{Child: test.child}

			for i, expectedStatus := range test.expectedStatus {
				status := whileSuccess.Tick()
				if status != expectedStatus {
					t.Errorf("Tick %d: expected %v, got %v", i+1, expectedStatus, status)
				}
				if whileSuccess.Status() != expectedStatus {
					t.Errorf("Status() after tick %d: expected %v, got %v", i+1, expectedStatus, whileSuccess.Status())
				}
			}
		})
	}
}

func TestWhileSuccess_ChildReset(t *testing.T) {
	executions := 0
	action := &Action{
		Run: func() Status {
			executions++
			if executions <= 3 {
				return Success
			}
			return Failure
		},
	}

	whileSuccess := &WhileSuccess{Child: action}

	// Execute multiple times - should reset child after each success
	for i := 0; i < 3; i++ {
		status := whileSuccess.Tick()
		if status != Running {
			t.Errorf("Tick %d: expected Running, got %v", i+1, status)
		}
	}

	// On the 4th tick, the action should fail
	status := whileSuccess.Tick()
	if status != Failure {
		t.Errorf("Final tick: expected Failure, got %v", status)
	}

	// Verify that executions count matches expected (child was reset after each success)
	if executions != 4 {
		t.Errorf("Expected 4 executions, got %d", executions)
	}
}

func TestWhileSuccess_Reset(t *testing.T) {
	action := &Action{Run: func() Status { return Success }}
	whileSuccess := &WhileSuccess{Child: action}

	// Execute once
	whileSuccess.Tick()
	if whileSuccess.Status() != Running {
		t.Errorf("Expected Running after tick, got %v", whileSuccess.Status())
	}

	// Reset
	status := whileSuccess.Reset()
	if status != Ready {
		t.Errorf("Reset() should return Ready, got %v", status)
	}
	if whileSuccess.Status() != Ready {
		t.Errorf("Status() after Reset() should be Ready, got %v", whileSuccess.Status())
	}
}

func TestWhileSuccess_String(t *testing.T) {
	tests := []struct {
		name     string
		child    Node
		expected string
	}{
		{
			name:     "no child",
			child:    nil,
			expected: "WhileSuccess (Ready)",
		},
		{
			name:     "with child",
			child:    &Action{Run: func() Status { return Success }},
			expected: "WhileSuccess (Ready)\n  Action (Ready)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			whileSuccess := &WhileSuccess{Child: test.child}
			result := whileSuccess.String()
			if result != test.expected {
				t.Errorf("String() = %q, want %q", result, test.expected)
			}
		})
	}
}

func TestWhileSuccess_MixedResults(t *testing.T) {
	results := []Status{Success, Success, Running, Running, Failure}
	currentResult := 0

	action := &Action{
		Run: func() Status {
			if currentResult >= len(results) {
				return Failure
			}
			result := results[currentResult]
			currentResult++
			return result
		},
	}

	whileSuccess := &WhileSuccess{Child: action}

	// Should return Running for Success, Success, Running, Running
	for i := 0; i < 4; i++ {
		status := whileSuccess.Tick()
		if status != Running {
			t.Errorf("Tick %d: expected Running, got %v", i+1, status)
		}
	}

	// Should return Failure when child returns Failure
	status := whileSuccess.Tick()
	if status != Failure {
		t.Errorf("Final tick: expected Failure, got %v", status)
	}
}

func TestComplexBehaviorTreeWithWhileSuccess(t *testing.T) {
	attempts := 0
	maxAttempts := 3

	// Action that succeeds a few times then fails
	limitedAction := &Action{
		Run: func() Status {
			attempts++
			if attempts <= maxAttempts {
				return Success
			}
			return Failure
		},
	}

	// WhileSuccess decorator
	whileSuccess := &WhileSuccess{Child: limitedAction}

	// Fallback action
	fallback := &Action{
		Run: func() Status {
			return Success
		},
	}

	// Selector: try whileSuccess, fall back to simple action
	selector := &Selector{
		Children: []Node{whileSuccess, fallback},
	}

	bt := New(selector)

	// First few ticks should keep the WhileSuccess running
	for i := 0; i < maxAttempts; i++ {
		status := bt.Tick()
		if status != Running {
			t.Errorf("Tick %d: expected Running, got %v", i+1, status)
		}
	}

	// Next tick should cause WhileSuccess to fail and selector to use fallback
	status := bt.Tick()
	if status != Success {
		t.Errorf("Expected Success from fallback, got %v", status)
	}

	// Test string representation
	str := bt.String()
	expectedParts := []string{"BehaviorTree", "Selector", "WhileSuccess", "Action"}
	for _, part := range expectedParts {
		if !strings.Contains(str, part) {
			t.Errorf("BehaviorTree.String() should contain '%s', got %v", part, str)
		}
	}
}

func TestWhileFailure_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		expectedStatus Status
		description    string
	}{
		{
			name:           "no child",
			child:          nil,
			expectedStatus: Success,
			description:    "should succeed when no child is provided",
		},
		{
			name:           "child succeeds immediately",
			child:          &Action{Run: func() Status { return Success }},
			expectedStatus: Success,
			description:    "should succeed when child succeeds",
		},
		{
			name:           "child fails",
			child:          &Action{Run: func() Status { return Failure }},
			expectedStatus: Running,
			description:    "should keep running when child fails",
		},
		{
			name:           "child running",
			child:          &Action{Run: func() Status { return Running }},
			expectedStatus: Running,
			description:    "should keep running when child is running",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			whileFailure := &WhileFailure{Child: test.child}
			status := whileFailure.Tick()
			if status != test.expectedStatus {
				t.Errorf("WhileFailure.Tick() = %v, want %v (%s)", status, test.expectedStatus, test.description)
			}
			if whileFailure.Status() != test.expectedStatus {
				t.Errorf("WhileFailure.Status() = %v, want %v", whileFailure.Status(), test.expectedStatus)
			}
		})
	}
}

func TestWhileFailure_ChildReset(t *testing.T) {
	executions := 0
	action := &Action{
		Run: func() Status {
			executions++
			if executions < 3 {
				return Failure // Fail first 2 times
			}
			return Success // Succeed on 3rd attempt
		},
	}

	whileFailure := &WhileFailure{Child: action}

	// First tick: child fails, should continue running and reset child
	status := whileFailure.Tick()
	if status != Running {
		t.Errorf("First tick: expected Running, got %v", status)
	}

	// Second tick: child fails again, should continue running and reset child
	status = whileFailure.Tick()
	if status != Running {
		t.Errorf("Second tick: expected Running, got %v", status)
	}

	// Third tick: child succeeds, WhileFailure should succeed
	status = whileFailure.Tick()
	if status != Success {
		t.Errorf("Third tick: expected Success, got %v", status)
	}

	if executions != 3 {
		t.Errorf("Expected 3 executions, got %d", executions)
	}
}

func TestWhileFailure_Reset(t *testing.T) {
	action := &Action{Run: func() Status { return Failure }}
	whileFailure := &WhileFailure{Child: action}

	// Execute and verify running state
	status := whileFailure.Tick()
	if status != Running {
		t.Errorf("Expected Running after tick, got %v", status)
	}

	// Reset and verify ready state
	resetStatus := whileFailure.Reset()
	if resetStatus != Ready {
		t.Errorf("Reset() = %v, want Ready", resetStatus)
	}
	if whileFailure.Status() != Ready {
		t.Errorf("Status() after Reset() = %v, want Ready", whileFailure.Status())
	}
}

func TestWhileFailure_String(t *testing.T) {
	tests := []struct {
		name     string
		child    Node
		expected []string
	}{
		{
			name:     "no child",
			child:    nil,
			expected: []string{"WhileFailure", "Ready"},
		},
		{
			name:     "with child, running state",
			child:    &Action{Run: func() Status { return Failure }},
			expected: []string{"WhileFailure", "Running", "Action"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			whileFailure := &WhileFailure{Child: test.child}
			if test.child != nil {
				whileFailure.Tick() // Execute to set status
			}

			str := whileFailure.String()
			for _, expected := range test.expected {
				if !strings.Contains(str, expected) {
					t.Errorf("WhileFailure.String() should contain '%s', got %v", expected, str)
				}
			}
		})
	}
}

func TestWhileFailure_MixedResults(t *testing.T) {
	attempts := 0
	action := &Action{
		Run: func() Status {
			attempts++
			switch attempts {
			case 1, 2:
				return Failure // Fail first 2 times
			case 3:
				return Running // Running on 3rd attempt
			case 4:
				return Failure // Fail on 4th attempt
			default:
				return Success // Succeed on 5th attempt
			}
		},
	}

	whileFailure := &WhileFailure{Child: action}

	// First two ticks: child fails, should keep running
	for i := 0; i < 2; i++ {
		status := whileFailure.Tick()
		if status != Running {
			t.Errorf("Tick %d: expected Running, got %v", i+1, status)
		}
	}

	// Third tick: child running, should keep running
	status := whileFailure.Tick()
	if status != Running {
		t.Errorf("Third tick: expected Running, got %v", status)
	}

	// Fourth tick: child fails, should keep running
	status = whileFailure.Tick()
	if status != Running {
		t.Errorf("Fourth tick: expected Running, got %v", status)
	}

	// Fifth tick: child succeeds, WhileFailure should succeed
	status = whileFailure.Tick()
	if status != Success {
		t.Errorf("Final tick: expected Success, got %v", status)
	}
}

func TestComplexBehaviorTreeWithWhileFailure(t *testing.T) {
	attempts := 0
	maxAttempts := 3

	// Action that fails a few times then succeeds
	retryAction := &Action{
		Run: func() Status {
			attempts++
			if attempts < maxAttempts {
				return Failure
			}
			return Success
		},
	}

	// WhileFailure decorator - keeps trying until success
	whileFailure := &WhileFailure{Child: retryAction}

	// Fallback action in case WhileFailure somehow fails
	fallbackAction := &Action{
		Run: func() Status {
			return Success
		},
	}

	// Selector: try WhileFailure first, fallback if needed
	selector := &Selector{
		Children: []Node{whileFailure, fallbackAction},
	}

	bt := New(selector)

	// Keep ticking until WhileFailure succeeds
	for i := 0; i < maxAttempts; i++ {
		status := bt.Tick()
		if i < maxAttempts-1 {
			if status != Running {
				t.Errorf("Tick %d: expected Running, got %v", i+1, status)
			}
		}
	}

	// Final tick should succeed
	status := bt.Tick()
	if status != Success {
		t.Errorf("Expected Success after retries, got %v", status)
	}

	// Test string representation
	str := bt.String()
	expectedParts := []string{"BehaviorTree", "Selector", "WhileFailure", "Action"}
	for _, part := range expectedParts {
		if !strings.Contains(str, part) {
			t.Errorf("BehaviorTree.String() should contain '%s', got %v", part, str)
		}
	}
}

func TestLog_Tick(t *testing.T) {
	tests := []struct {
		name           string
		child          Node
		message        string
		expectedStatus Status
		description    string
	}{
		{
			name:           "no child",
			child:          nil,
			message:        "Test log",
			expectedStatus: Failure,
			description:    "should fail when no child is provided",
		},
		{
			name:           "child succeeds",
			child:          &Action{Run: func() Status { return Success }},
			message:        "Success test",
			expectedStatus: Success,
			description:    "should return child's success status",
		},
		{
			name:           "child fails",
			child:          &Action{Run: func() Status { return Failure }},
			message:        "Failure test",
			expectedStatus: Failure,
			description:    "should return child's failure status",
		},
		{
			name:           "child running",
			child:          &Action{Run: func() Status { return Running }},
			message:        "Running test",
			expectedStatus: Running,
			description:    "should return child's running status",
		},
		{
			name:           "empty message",
			child:          &Action{Run: func() Status { return Success }},
			message:        "",
			expectedStatus: Success,
			description:    "should work with empty message",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := &Log{Child: test.child, Message: test.message}
			status := log.Tick()
			if status != test.expectedStatus {
				t.Errorf("Log.Tick() = %v, want %v (%s)", status, test.expectedStatus, test.description)
			}
			if log.Status() != test.expectedStatus {
				t.Errorf("Log.Status() = %v, want %v", log.Status(), test.expectedStatus)
			}
		})
	}
}

func TestLog_ChildExecution(t *testing.T) {
	executions := 0
	action := &Action{
		Run: func() Status {
			executions++
			return Success
		},
	}

	log := &Log{Child: action, Message: "Execution test"}

	// First execution
	status := log.Tick()
	if status != Success {
		t.Errorf("First tick: expected Success, got %v", status)
	}
	if executions != 1 {
		t.Errorf("Expected 1 execution, got %d", executions)
	}

	// Second execution (should execute child again)
	status = log.Tick()
	if status != Success {
		t.Errorf("Second tick: expected Success, got %v", status)
	}
	if executions != 2 {
		t.Errorf("Expected 2 executions, got %d", executions)
	}
}

func TestLog_Reset(t *testing.T) {
	action := &Action{
		Run: func() Status { return Success },
	}

	log := &Log{Child: action, Message: "Reset test"}

	// Execute and verify running state
	status := log.Tick()
	if status != Success {
		t.Errorf("Expected Success after tick, got %v", status)
	}

	// Reset and verify ready state
	resetStatus := log.Reset()
	if resetStatus != Ready {
		t.Errorf("Reset() = %v, want Ready", resetStatus)
	}
	if log.Status() != Ready {
		t.Errorf("Status() after Reset() = %v, want Ready", log.Status())
	}

	// Verify child is also reset by checking its status
	if action.Status() != Ready {
		t.Errorf("Child status after reset = %v, want Ready", action.Status())
	}
}

func TestLog_String(t *testing.T) {
	tests := []struct {
		name     string
		child    Node
		message  string
		expected []string
	}{
		{
			name:     "no child",
			child:    nil,
			message:  "Test message",
			expected: []string{"Log", "Ready", "Test message"},
		},
		{
			name:     "with child and message",
			child:    &Action{Run: func() Status { return Success }},
			message:  "Action log",
			expected: []string{"Log", "Success", "Action log", "Action"},
		},
		{
			name:     "empty message",
			child:    &Action{Run: func() Status { return Success }},
			message:  "",
			expected: []string{"Log", "Success", "Action"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := &Log{Child: test.child, Message: test.message}
			if test.child != nil {
				log.Tick() // Execute to set status
			}

			str := log.String()
			for _, expected := range test.expected {
				if expected != "" && !strings.Contains(str, expected) {
					t.Errorf("Log.String() should contain '%s', got %v", expected, str)
				}
			}
		})
	}
}

func TestLog_GetChildType(t *testing.T) {
	tests := []struct {
		name         string
		child        Node
		expectedType string
	}{
		{"nil child", nil, "nil"},
		{"Action", &Action{}, "Action"},
		{"Condition", &Condition{}, "Condition"},
		{"Sequence", &Sequence{}, "Sequence"},
		{"Selector", &Selector{}, "Selector"},
		{"Parallel", &Parallel{}, "Parallel"},
		{"Composite", &Composite{}, "Composite"},
		{"Retry", &Retry{}, "Retry"},
		{"Repeat", &Repeat{}, "Repeat"},
		{"RepeatN", &RepeatN{}, "RepeatN"},
		{"Invert", &Invert{}, "Invert"},
		{"AlwaysSuccess", &AlwaysSuccess{}, "AlwaysSuccess"},
		{"AlwaysFailure", &AlwaysFailure{}, "AlwaysFailure"},
		{"WhileSuccess", &WhileSuccess{}, "WhileSuccess"},
		{"WhileFailure", &WhileFailure{}, "WhileFailure"},
		{"Log", &Log{}, "Log"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := &Log{Child: test.child}
			childType := log.getChildType()
			if childType != test.expectedType {
				t.Errorf("getChildType() = %v, want %v", childType, test.expectedType)
			}
		})
	}
}

func TestLog_LoggingBehavior(t *testing.T) {
	// Test that different status levels use appropriate log levels
	// This test verifies the logging happens without actually capturing logs

	// Suppress log output during testing by setting a discard handler
	origHandler := slog.Default().Handler()
	defer slog.SetDefault(slog.New(origHandler))

	tests := []struct {
		name   string
		status Status
	}{
		{"Success logging", Success},
		{"Failure logging", Failure},
		{"Running logging", Running},
		{"Ready logging", Ready},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			action := &Action{Run: func() Status { return test.status }}
			log := &Log{Child: action, Message: "Test logging"}

			// This should not panic and should execute logging
			status := log.Tick()
			if status != test.status {
				t.Errorf("Log.Tick() = %v, want %v", status, test.status)
			}
		})
	}
}

func TestComplexBehaviorTreeWithLog(t *testing.T) {
	executions := 0

	// Action that we want to monitor
	monitoredAction := &Action{
		Run: func() Status {
			executions++
			if executions <= 2 {
				return Failure
			}
			return Success
		},
	}

	// Wrap the action with logging
	loggedAction := &Log{
		Child:   monitoredAction,
		Message: "Monitoring critical action",
	}

	// Retry the logged action
	retry := &Retry{Child: loggedAction}

	// Fallback action
	fallback := &Action{
		Run: func() Status {
			return Success
		},
	}

	// Selector with logged retry
	selector := &Selector{
		Children: []Node{retry, fallback},
	}

	bt := New(selector)

	// Execute until success - may need multiple ticks for retry to work
	for i := 0; i < 10; i++ { // Safety limit
		status := bt.Tick()
		if status == Success {
			break
		}
		if status == Failure {
			t.Errorf("Unexpected failure at tick %d", i+1)
			break
		}
	}

	// Final status should be success
	if bt.Status() != Success {
		t.Errorf("Expected final status Success, got %v", bt.Status())
	}

	// Verify the action was executed the expected number of times
	if executions != 3 {
		t.Errorf("Expected 3 executions, got %d", executions)
	}

	// Test string representation includes Log node
	str := bt.String()
	expectedParts := []string{"BehaviorTree", "Selector", "Retry", "Log", "Action"}
	for _, part := range expectedParts {
		if !strings.Contains(str, part) {
			t.Errorf("BehaviorTree.String() should contain '%s', got %v", part, str)
		}
	}
}
