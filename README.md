
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
- **Retry**: Decorator node. Retries its child until it succeeds, ignoring all failures. Returns Success when child succeeds, Running while retrying.
- **Repeat**: Decorator node. Repeats its child until it fails. Returns Running while child succeeds (and resets it), Failure when child fails.

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

// Retry until success (keeps trying after failures)
retry := &behave.Retry{
    Child: &behave.Action{Run: unreliableFunc},
}

tree := behave.New(par)
```

## Retry Node Example

The Retry node is particularly useful for unreliable operations that might fail but should eventually succeed:

```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/rbrabson/behave"
)

func main() {
    attempts := 0
    
    // Unreliable action that fails randomly
    unreliableAction := &behave.Action{
        Run: func() behave.Status {
            attempts++
            fmt.Printf("Attempt %d: ", attempts)
            
            if rand.Float32() < 0.7 { // 70% chance of failure
                fmt.Println("Failed!")
                return behave.Failure
            }
            
            fmt.Println("Success!")
            return behave.Success
        },
    }
    
    // Wrap with Retry to keep trying until success
    retry := &behave.Retry{Child: unreliableAction}
    tree := behave.New(retry)
    
    // Keep ticking until the tree succeeds
    for tree.Status() != behave.Success {
        tree.Tick()
    }
    
    fmt.Printf("Finally succeeded after %d attempts!\n", attempts)
}
```

## Repeat Node Example

The Repeat node is useful for tasks that should continue running until they fail:

```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/rbrabson/behave"
)

func main() {
    rounds := 0
    
    // Action that succeeds a few times then fails
    taskAction := &behave.Action{
        Run: func() behave.Status {
            rounds++
            fmt.Printf("Round %d: ", rounds)
            
            if rounds >= 5 { // Fail after 5 rounds
                fmt.Println("Failed - stopping!")
                return behave.Failure
            }
            
            fmt.Println("Success - continuing!")
            return behave.Success
        },
    }
    
    // Wrap with Repeat to keep running until failure
    repeat := &behave.Repeat{Child: taskAction}
    tree := behave.New(repeat)
    
    // Keep ticking until the tree fails
    for tree.Status() != behave.Failure {
        status := tree.Tick()
        fmt.Printf("Tree status: %s\n", status.String())
    }
    
    fmt.Printf("Stopped after %d rounds!\n", rounds)
}
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
