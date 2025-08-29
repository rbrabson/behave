
# Behavior Tree

## Overview

This package provides a simple and extensible implementation of a **Behavior Tree** in Go. Behavior trees are widely used in AI for games and robotics to model complex, hierarchical decision-making logic in a modular and reusable way.

## Core Concepts

### Status

Nodes in the tree return a `Status` after execution:

- `Success`: Node completed successfully
- `Failure`: Node failed
- `Ready`: Node is ready to run
- `Running`: Node is currently running

### Node Interface

All nodes implement the following interface:

```go
type Node interface {
    Tick() Status   // Run the node on each tick
    Reset() Status  // Reset the node to initial state
    Status() Status // Get the current status
    String() string // String representation
}
```

### Node Types

- **Action**: Leaf node that performs an action. You provide a `Run` function.
- **Condition**: Leaf node that checks a condition. You provide a `Check` function.
- **Composite**: Combines a condition with any other node. First checks the condition, and if it succeeds, runs the child node.
- **Sequence**: Composite node. Runs children in order; fails or returns running if any child fails or is running, succeeds if all succeed.
- **Selector**: Composite node. Runs children in order; succeeds or returns running if any child succeeds or is running, fails if all fail.
- **Parallel**: Composite node. Runs all children in parallel; succeeds if at least `MinSuccessCount` children succeed.

### BehaviorTree

The `BehaviorTree` struct manages the root node and provides methods to tick, reset, and get the status of the tree.

## Example Usage

```go
package main

import (
    "fmt"
    "github.com/rbrabson/behave"
)

func main() {
    action := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Action running!")
            return behave.Success
        },
    }
    tree := behave.New(action)
    status := tree.Tick()
    fmt.Println("Tree status:", status)
}
```

## Tree Structure Example

You can compose trees using different node types:

```go
// Simple sequence
seq := &behave.Sequence{
    Children: []behave.Node{
        &behave.Condition{Check: myCheckFunc},
        &behave.Action{Run: myActionFunc},
    },
}

// Selector with fallback
sel := &behave.Selector{
    Children: []behave.Node{seq, &behave.Action{Run: fallbackFunc}},
}

// Composite (condition + action)
comp := &behave.Composite{
    Condition: &behave.Condition{Check: guardFunc},
    Child: &behave.Action{Run: protectedFunc},
}

// Parallel execution (at least 2 out of 3 must succeed)
par := &behave.Parallel{
    MinSuccessCount: 2,
    Children: []behave.Node{
        &behave.Action{Run: task1},
        &behave.Action{Run: task2},
        &behave.Action{Run: task3},
    },
}

tree := behave.New(par)
```

## Reset Functionality

All nodes can be reset to their initial state:

```go
// Reset the entire tree
tree.Reset()

// Reset individual nodes
node.Reset()
```

## Status String Representation

Each node and the tree itself can be printed for debugging:

```go
fmt.Println(tree.String())
```

## License

MIT
