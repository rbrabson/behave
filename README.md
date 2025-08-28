
# Behavior Tree

## Overview

This package provides a simple and extensible implementation of a **Behavior Tree** in Go. Behavior trees are widely used in AI for games and robotics to model complex, hierarchical decision-making logic in a modular and reusable way.

## Core Concepts

### Status

Nodes in the tree return a `Status` after execution:

- `Success`: Node completed successfully
- `Failure`: Node failed
- `Initializing`, `Ready`, `Running`, `Stopping`, `Stopped`: Used for lifecycle management

### Node Interface

All nodes implement the following interface:

```go
type Node interface {
    Init() Status   // Initialize the node
    Tick() Status   // Run the node on each tick
    Stop() Status   // Stop the node
    Status() Status // Get the current status
    String() string // String representation
}
```

### Node Types

- **Action**: Leaf node that performs an action. You provide `InitFunc`, `Run`, and `StopFunc` functions.
- **Condition**: Leaf node that checks a condition. You provide `InitFunc`, `Check`, and `StopFunc` functions.
- **Sequence**: Composite node. Runs children in order; fails or returns running if any child fails or is running, succeeds if all succeed.
- **Selector**: Composite node. Runs children in order; succeeds or returns running if any child succeeds or is running, fails if all fail.

### BehaviorTree

The `BehaviorTree` struct manages the root node and provides methods to initialize, tick, stop, and get the status of the tree.

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
    tree.Init()
    status := tree.Tick()
    fmt.Println("Tree status:", status)
}
```

## Tree Structure Example

You can compose trees using `Sequence` and `Selector` nodes:

```go
tree := behaviortree.New(sel)
seq := &behave.Sequence{
    Children: []behave.Node{
        &behave.Condition{Check: myCheckFunc},
        &behave.Action{Run: myActionFunc},
    },
}
sel := &behave.Selector{
    Children: []behave.Node{seq, otherNode},
}
tree := behave.New(sel)
```

## Status String Representation

Each node and the tree itself can be printed for debugging:

```go
fmt.Println(tree.String())
```

## License

MIT
