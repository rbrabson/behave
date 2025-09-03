
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

#### Leaf Nodes

- **Action**: Performs an action. You provide a `Run` function.
- **Condition**: Checks a condition. You provide a `Check` function.

#### Composite Nodes

- **Composite**: Combines a condition with any other node. First checks the condition, and if it succeeds, runs the child node.
- **Sequence**: Runs children in order; fails or returns running if any child fails or is running, succeeds if all succeed.
- **Selector**: Runs children in order; succeeds or returns running if any child succeeds or is running, fails if all fail.
- **Parallel**: Runs all children in parallel; succeeds if at least `MinSuccessCount` children succeed.

#### Decorator Nodes

- **Retry**: Retries its child until it succeeds, ignoring all failures. Returns Success when child succeeds, Running while retrying.
- **Repeat**: Repeats its child until it fails. Returns Running while child succeeds (and resets it), Failure when child fails.
- **RepeatN**: Executes its child a specific number of times before returning the child's last result. Returns Running until MaxCount is reached.
- **WhileSuccess**: Returns Running as long as its child is either Running or Success, and returns Failure otherwise. Useful for creating loops.
- **WhileFailure**: Returns Running as long as its child is either Running or Failure, and returns Success otherwise. Useful for retry loops.
- **Invert**: Inverts the result of its child. Changes Success to Failure and Failure to Success. Running and Ready states pass through unchanged.
- **AlwaysSuccess**: Always returns Success regardless of its child's result. Useful for ensuring a branch always succeeds.
- **AlwaysFailure**: Always returns Failure regardless of its child's result. Useful for ensuring a branch always fails.

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

// Repeat a specific number of times
repeatN := &behave.RepeatN{
    MaxCount: 3,
    Child: &behave.Action{Run: limitedFunc},
}

// Keep running while child succeeds or is running
whileSuccess := &behave.WhileSuccess{
    Child: &behave.Condition{Check: keepGoingFunc},
}

// Keep running while child fails or is running (retry pattern)
whileFailure := &behave.WhileFailure{
    Child: &behave.Action{Run: retryableFunc},
}

// Invert a condition (succeeds when condition fails)
invertedCondition := &behave.Invert{
    Child: &behave.Condition{Check: avoidThisFunc},
}

// Always succeed (useful for optional tasks)
alwaysSuccess := &behave.AlwaysSuccess{
    Child: &behave.Action{Run: optionalFunc},
}

// Always fail (useful for testing or negative conditions)
alwaysFailure := &behave.AlwaysFailure{
    Child: &behave.Action{Run: shouldFailFunc},
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

## Invert Node Example

The Invert node is useful for negating conditions or creating "avoid" behaviors:

```go
package main

import (
    "fmt"
    "github.com/rbrabson/behave"
)

func main() {
    isEnemyNear := false
    
    // Condition that checks if enemy is near
    enemyCheck := &behave.Condition{
        Check: func() behave.Status {
            fmt.Printf("Enemy near: %t\n", isEnemyNear)
            if isEnemyNear {
                return behave.Success
            }
            return behave.Failure
        },
    }
    
    // Invert the condition - succeeds when NO enemy is near
    notEnemyNear := &behave.Invert{Child: enemyCheck}
    
    // Action to perform when safe
    safeAction := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Performing safe action!")
            return behave.Success
        },
    }
    
    // Sequence: only do safe action when no enemy is near
    sequence := &behave.Sequence{
        Children: []behave.Node{notEnemyNear, safeAction},
    }
    
    tree := behave.New(sequence)
    
    // Test with no enemy
    fmt.Println("=== No enemy present ===")
    status := tree.Tick()
    fmt.Printf("Result: %s\n\n", status.String())
    
    // Test with enemy present
    isEnemyNear = true
    tree.Reset()
    fmt.Println("=== Enemy detected ===")
    status = tree.Tick()
    fmt.Printf("Result: %s\n", status.String())
}
```

## RepeatN Node Example

The RepeatN node executes its child a specific number of times, which is useful for controlled repetition:

```go
package main

import (
    "fmt"
    "github.com/rbrabson/behave"
)

func main() {
    executions := 0
    
    // Action that counts executions
    countAction := &behave.Action{
        Run: func() behave.Status {
            executions++
            fmt.Printf("Execution %d\n", executions)
            return behave.Success
        },
    }
    
    // Repeat exactly 3 times
    repeatN := &behave.RepeatN{
        MaxCount: 3,
        Child: countAction,
    }
    
    tree := behave.New(repeatN)
    
    // Keep ticking until completed
    for tree.Status() != behave.Success {
        status := tree.Tick()
        fmt.Printf("Tree status: %s\n", status.String())
    }
    
    fmt.Printf("Completed after %d executions!\n", executions)
}
```

## AlwaysSuccess and AlwaysFailure Node Examples

These decorator nodes are useful for ensuring specific outcomes regardless of their child's result:

```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/rbrabson/behave"
)

func main() {
    // Unreliable action that might fail
    unreliableAction := &behave.Action{
        Run: func() behave.Status {
            if rand.Float32() < 0.5 {
                fmt.Println("Action failed!")
                return behave.Failure
            }
            fmt.Println("Action succeeded!")
            return behave.Success
        },
    }
    
    // Wrap with AlwaysSuccess to ensure it never fails the tree
    optionalTask := &behave.AlwaysSuccess{Child: unreliableAction}
    
    // Critical action that must execute
    criticalAction := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Critical action executed!")
            return behave.Success
        },
    }
    
    // Sequence that includes the optional task
    sequence := &behave.Sequence{
        Children: []behave.Node{optionalTask, criticalAction},
    }
    
    tree := behave.New(sequence)
    
    fmt.Println("=== Running sequence with optional task ===")
    status := tree.Tick()
    fmt.Printf("Final result: %s\n\n", status.String())
    
    // Example with AlwaysFailure for testing
    failingCondition := &behave.AlwaysFailure{
        Child: &behave.Action{Run: func() behave.Status { return behave.Success }},
    }
    
    fallbackAction := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Fallback action executed!")
            return behave.Success
        },
    }
    
    // Selector that always goes to fallback
    selector := &behave.Selector{
        Children: []behave.Node{failingCondition, fallbackAction},
    }
    
    fallbackTree := behave.New(selector)
    
    fmt.Println("=== Running selector with always-failing condition ===")
    status = fallbackTree.Tick()
    fmt.Printf("Final result: %s\n", status.String())
}
```

## WhileSuccess Node Example

The WhileSuccess node creates loops that continue while a condition remains true or while an action keeps succeeding:

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
    "github.com/rbrabson/behave"
)

func main() {
    attempts := 0
    maxAttempts := 5
    
    // Condition that succeeds a limited number of times
    resourceCheck := &behave.Condition{
        Check: func() behave.Status {
            attempts++
            fmt.Printf("Attempt %d: ", attempts)
            
            if attempts <= maxAttempts {
                fmt.Println("Resources available - continuing!")
                return behave.Success
            }
            
            fmt.Println("Resources depleted - stopping!")
            return behave.Failure
        },
    }
    
    // WhileSuccess will keep running as long as resources are available
    whileSuccess := &behave.WhileSuccess{Child: resourceCheck}
    
    tree := behave.New(whileSuccess)
    
    fmt.Println("=== WhileSuccess Loop Example ===")
    
    // Keep ticking until the condition fails
    for tree.Status() != behave.Failure {
        status := tree.Tick()
        fmt.Printf("Tree status: %s\n", status.String())
        
        // Add a small delay to see the loop in action
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Printf("Loop completed after %d attempts!\n\n", attempts)
    
    // Example with an action that sometimes fails
    taskAttempts := 0
    
    unreliableTask := &behave.Action{
        Run: func() behave.Status {
            taskAttempts++
            fmt.Printf("Task attempt %d: ", taskAttempts)
            
            // 70% chance of success
            if rand.Float32() < 0.7 {
                fmt.Println("Task succeeded!")
                return behave.Success
            }
            
            fmt.Println("Task failed - stopping loop!")
            return behave.Failure
        },
    }
    
    whileTaskSuccess := &behave.WhileSuccess{Child: unreliableTask}
    taskTree := behave.New(whileTaskSuccess)
    
    fmt.Println("=== WhileSuccess with Unreliable Task ===")
    
    // Keep ticking until the task fails
    for taskTree.Status() != behave.Failure {
        status := taskTree.Tick()
        fmt.Printf("Task tree status: %s\n", status.String())
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Printf("Task loop completed after %d attempts!\n", taskAttempts)
}
```

## WhileFailure Node Example

The WhileFailure node creates retry loops that continue while a task fails, and succeeds when the task finally succeeds:

```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/rbrabson/behave"
)

func main() {
    attempts := 0
    maxAttempts := 5

    // Unreliable network operation that might fail
    networkOperation := &behave.Action{
        Run: func() behave.Status {
            attempts++
            fmt.Printf("Network attempt %d: ", attempts)
            
            // Simulate increasing success probability over time
            successChance := float32(attempts) / float32(maxAttempts)
            if rand.Float32() < successChance {
                fmt.Println("Success!")
                return behave.Success
            }
            
            fmt.Println("Failed - retrying...")
            return behave.Failure
        },
    }
    
    // WhileFailure will keep trying until the operation succeeds
    retryLoop := &behave.WhileFailure{Child: networkOperation}
    
    tree := behave.New(retryLoop)
    
    fmt.Println("=== WhileFailure Retry Loop Example ===")
    
    // Keep ticking until success
    for tree.Status() != behave.Success {
        status := tree.Tick()
        fmt.Printf("Tree status: %s\n", status.String())
        if attempts >= maxAttempts {
            break // Safety break
        }
    }
    
    fmt.Printf("Operation completed after %d attempts!\n", attempts)
    
    // Example 2: WhileFailure in a larger behavior tree
    fmt.Println("\n=== WhileFailure in Complex Tree ===")
    
    connectionAttempts := 0
    connectAction := &behave.Action{
        Run: func() behave.Status {
            connectionAttempts++
            fmt.Printf("Connection attempt %d: ", connectionAttempts)
            
            if connectionAttempts >= 3 {
                fmt.Println("Connected!")
                return behave.Success
            }
            
            fmt.Println("Connection failed")
            return behave.Failure
        },
    }
    
    // Task to execute after successful connection
    mainTask := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Executing main task...")
            return behave.Success
        },
    }
    
    // WhileFailure for connection retry
    connectionRetry := &behave.WhileFailure{Child: connectAction}
    
    // Sequence: connect then execute task
    sequence := &behave.Sequence{
        Children: []behave.Node{connectionRetry, mainTask},
    }
    
    complexTree := behave.New(sequence)
    
    // Execute the complex tree
    for complexTree.Status() != behave.Success {
        status := complexTree.Tick()
        fmt.Printf("Complex tree status: %s\n", status.String())
    }
    
    fmt.Printf("Complex operation completed after %d connection attempts!\n", connectionAttempts)
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
