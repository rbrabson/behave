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
	var action *Action
	action = &Action{
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
