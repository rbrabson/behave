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
		{Initializing, "Initializing"},
		{Ready, "Ready"},
		{Running, "Running"},
		{Stopping, "Stopping"},
		{Stopped, "Stopped"},
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
	action := &Action{
		InitFunc: func() Status { return Ready },
	}
	bt := New(action)

	status := bt.Init()
	if status != Ready {
		t.Errorf("BehaviorTree.Init() = %v, want %v", status, Ready)
	}
	if bt.Status() != Ready {
		t.Errorf("BehaviorTree.Status() after Init() = %v, want %v", bt.Status(), Ready)
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

func TestBehaviorTree_Stop(t *testing.T) {
	action := &Action{
		StopFunc: func() Status { return Stopped },
	}
	bt := New(action)

	status := bt.Stop()
	if status != Stopped {
		t.Errorf("BehaviorTree.Stop() = %v, want %v", status, Stopped)
	}
	if bt.Status() != Stopped {
		t.Errorf("BehaviorTree.Status() after Stop() = %v, want %v", bt.Status(), Stopped)
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

func TestAction_Init(t *testing.T) {
	tests := []struct {
		name     string
		action   *Action
		expected Status
	}{
		{
			name:     "no init func",
			action:   &Action{},
			expected: Ready,
		},
		{
			name: "with init func",
			action: &Action{
				InitFunc: func() Status { return Running },
			},
			expected: Running,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := test.action.Init()
			if status != test.expected {
				t.Errorf("Action.Init() = %v, want %v", status, test.expected)
			}
			if test.action.Status() != test.expected {
				t.Errorf("Action.Status() after Init() = %v, want %v", test.action.Status(), test.expected)
			}
		})
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

func TestAction_Stop(t *testing.T) {
	tests := []struct {
		name     string
		action   *Action
		expected Status
	}{
		{
			name:     "no stop func",
			action:   &Action{},
			expected: Stopped,
		},
		{
			name: "with stop func",
			action: &Action{
				StopFunc: func() Status { return Stopping },
			},
			expected: Stopping,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status := test.action.Stop()
			if status != test.expected {
				t.Errorf("Action.Stop() = %v, want %v", status, test.expected)
			}
			if test.action.Status() != test.expected {
				t.Errorf("Action.Status() after Stop() = %v, want %v", test.action.Status(), test.expected)
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

	// Initialize the tree
	bt.Init()
	if bt.Status() != Ready {
		t.Errorf("BehaviorTree.Status() after Init() = %v, want %v", bt.Status(), Ready)
	}

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
