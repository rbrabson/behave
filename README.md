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
- **Parallel**: Runs all children in parallel; succeeds if at least `MinSuccessCount` children succeed, fails if it becomes impossible to reach MinSuccessCount (too many failures), and returns Running while children are still executing.

#### Decorator Nodes

- **Retry**: Retries its child until it succeeds, ignoring all failures. Returns Success when child succeeds, Running while retrying.
- **Repeat**: Repeats its child node until the child returns Failure. Returns Running while the child returns Success or Running, and returns Failure when the child fails. Useful for tasks that should continue until a failure occurs.
- **RepeatN**: Executes its child a specific number of times (MaxCount). Returns Running while the execution count is below MaxCount, then returns the child's final result. Useful for controlled repetition. See example below.
- **Forever**: Runs its child forever, always returning Running and ignoring the child's status. Useful for infinite loops or background tasks.
- **WhileSuccess**: Repeatedly runs its child as long as it returns Success or Running. Returns Running while the child succeeds (resetting it for the next iteration) or is running, and returns Failure when the child fails. Useful for creating loops that continue until failure.
- **WhileFailure**: Repeatedly runs its child as long as it returns Failure or Running. Returns Running while the child fails (resetting it for retry) or is running, and returns Success when the child succeeds. Useful for retry loops that continue until success.
- **Invert**: Inverts the result of its child. Changes Success to Failure and Failure to Success. Running and Ready states pass through unchanged.
- **AlwaysSuccess**: Always returns Success regardless of its child's result. Useful for ensuring a branch always succeeds.
- **AlwaysFailure**: Always returns Failure regardless of its child's result. Useful for ensuring a branch always fails.
- **Log**: Executes its child and logs the result using structured logging (slog). Returns the child's status unchanged. Supports custom log levels or uses defaults (Info for Success, Warn for Failure, Debug for Running/Ready). Useful for debugging and monitoring.

- **WithTimeout**: Runs its child node for at most the specified duration (using Go's `time.Duration`). If the child completes (returns Success or Failure) before the duration expires, WithTimeout returns that status immediately. If the duration expires while the child is still running (status == Ready or Running), WithTimeout returns Failure. Useful for time-limited behaviors, polling, or enforcing timeouts.

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

## Action Node Example

The Action node is the most basic leaf node that performs an operation:

```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/rbrabson/behave"
)

func main() {
    attempts := 0
    
    // Simple action that always succeeds
    simpleAction := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Simple action executed!")
            return behave.Success
        },
    }
    
    // Action that might fail
    unreliableAction := &behave.Action{
        Run: func() behave.Status {
            attempts++
            fmt.Printf("Unreliable action attempt %d: ", attempts)
            
            if rand.Float32() < 0.3 { // 30% success rate
                fmt.Println("Success!")
                return behave.Success
            }
            
            fmt.Println("Failed!")
            return behave.Failure
        },
    }
    
    // Long-running action
    ticks := 0
    longRunningAction := &behave.Action{
        Run: func() behave.Status {
            ticks++
            fmt.Printf("Long running action tick %d: ", ticks)
            
            if ticks < 3 {
                fmt.Println("Still running...")
                return behave.Running
            }
            
            fmt.Println("Completed!")
            return behave.Success
        },
    }
    
    // Test simple action
    tree1 := behave.New(simpleAction)
    fmt.Println("=== Simple Action ===")
    fmt.Printf("Result: %s\n\n", tree1.Tick().String())
    
    // Test unreliable action
    tree2 := behave.New(unreliableAction)
    fmt.Println("=== Unreliable Action ===")
    for tree2.Status() != behave.Success && attempts < 5 {
        tree2.Tick()
        tree2.Reset() // Reset for next attempt
    }
    fmt.Printf("Final result: %s\n\n", tree2.Status().String())
    
    // Test long-running action
    tree3 := behave.New(longRunningAction)
    fmt.Println("=== Long Running Action ===")
    for tree3.Status() != behave.Success {
        status := tree3.Tick()
        fmt.Printf("Status: %s\n", status.String())
    }
}
```

## Condition Node Example

The Condition node checks a condition and returns Success or Failure (never Running):

```go
package main

import (
    "fmt"
    "time"
    "github.com/rbrabson/behave"
)

func main() {
    // System state
    temperature := 20.0
    isSystemReady := false
    
    // Simple condition check
    temperatureCheck := &behave.Condition{
        Check: func() behave.Status {
            fmt.Printf("Checking temperature: %.1f°C\n", temperature)
            if temperature < 25.0 {
                fmt.Println("Temperature OK")
                return behave.Success
            }
            fmt.Println("Temperature too high!")
            return behave.Failure
        },
    }
    
    // System readiness check
    readinessCheck := &behave.Condition{
        Check: func() behave.Status {
            fmt.Printf("System ready: %t\n", isSystemReady)
            if isSystemReady {
                return behave.Success
            }
            return behave.Failure
        },
    }
    
    // Complex condition with multiple factors
    complexCheck := &behave.Condition{
        Check: func() behave.Status {
            hour := time.Now().Hour()
            fmt.Printf("Current hour: %d\n", hour)
            
            // Business hours check (9 AM to 5 PM)
            if hour >= 9 && hour < 17 {
                fmt.Println("Within business hours")
                return behave.Success
            }
            
            fmt.Println("Outside business hours")
            return behave.Failure
        },
    }
    
    // Test temperature condition
    fmt.Println("=== Temperature Check ===")
    tree1 := behave.New(temperatureCheck)
    fmt.Printf("Result: %s\n\n", tree1.Tick().String())
    
    // Change temperature and test again
    temperature = 30.0
    tree1.Reset()
    fmt.Printf("After temperature change: %s\n\n", tree1.Tick().String())
    
    // Test system readiness
    fmt.Println("=== System Readiness Check ===")
    tree2 := behave.New(readinessCheck)
    fmt.Printf("Before system ready: %s\n", tree2.Tick().String())
    
    isSystemReady = true
    tree2.Reset()
    fmt.Printf("After system ready: %s\n\n", tree2.Tick().String())
    
    // Test complex condition
    fmt.Println("=== Business Hours Check ===")
    tree3 := behave.New(complexCheck)
    fmt.Printf("Result: %s\n", tree3.Tick().String())
}
```

## Sequence Node Example

The Sequence node executes children in order, failing if any child fails:

```go
package main

import (
    "fmt"
    "github.com/rbrabson/behave"
)

func main() {
    step := 0
    
    // Create sequence steps
    step1 := &behave.Action{
        Run: func() behave.Status {
            step++
            fmt.Printf("Step 1: Initialize (step=%d)\n", step)
            return behave.Success
        },
    }
    
    step2 := &behave.Action{
        Run: func() behave.Status {
            step++
            fmt.Printf("Step 2: Process (step=%d)\n", step)
            return behave.Success
        },
    }
    
    step3 := &behave.Action{
        Run: func() behave.Status {
            step++
            fmt.Printf("Step 3: Finalize (step=%d)\n", step)
            return behave.Success
        },
    }
    
    // Successful sequence
    successSequence := &behave.Sequence{
        Children: []behave.Node{step1, step2, step3},
    }
    
    fmt.Println("=== Successful Sequence ===")
    tree1 := behave.New(successSequence)
    status := tree1.Tick()
    fmt.Printf("Final result: %s\n\n", status.String())
    
    // Sequence with failure
    step = 0
    failingStep := &behave.Action{
        Run: func() behave.Status {
            step++
            fmt.Printf("Failing step (step=%d)\n", step)
            return behave.Failure
        },
    }
    
    remainingStep := &behave.Action{
        Run: func() behave.Status {
            step++
            fmt.Printf("This should not execute (step=%d)\n", step)
            return behave.Success
        },
    }
    
    failureSequence := &behave.Sequence{
        Children: []behave.Node{step1, failingStep, remainingStep},
    }
    
    fmt.Println("=== Sequence with Failure ===")
    tree2 := behave.New(failureSequence)
    status = tree2.Tick()
    fmt.Printf("Final result: %s (step=%d)\n\n", status.String(), step)
    
    // Sequence with running state
    step = 0
    runningTicks := 0
    runningStep := &behave.Action{
        Run: func() behave.Status {
            runningTicks++
            fmt.Printf("Running step tick %d\n", runningTicks)
            if runningTicks < 3 {
                return behave.Running
            }
            return behave.Success
        },
    }
    
    runningSequence := &behave.Sequence{
        Children: []behave.Node{step1, runningStep, step3},
    }
    
    fmt.Println("=== Sequence with Running State ===")
    tree3 := behave.New(runningSequence)
    for tree3.Status() != behave.Success {
        status = tree3.Tick()
        fmt.Printf("Status: %s\n", status.String())
    }
}
```

## Selector Node Example

The Selector node tries children in order until one succeeds:

```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/rbrabson/behave"
)

func main() {
    attempt := 0
    
    // Primary option (might fail)
    primaryOption := &behave.Action{
        Run: func() behave.Status {
            attempt++
            fmt.Printf("Primary option attempt %d: ", attempt)
            
            if rand.Float32() < 0.3 { // 30% success rate
                fmt.Println("Success!")
                return behave.Success
            }
            
            fmt.Println("Failed!")
            return behave.Failure
        },
    }
    
    // Secondary option (more reliable)
    secondaryOption := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Secondary option: Success!")
            return behave.Success
        },
    }
    
    // Fallback option (always works)
    fallbackOption := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Fallback option: Success!")
            return behave.Success
        },
    }
    
    // Create selector with fallback chain
    selector := &behave.Selector{
        Children: []behave.Node{primaryOption, secondaryOption, fallbackOption},
    }
    
    fmt.Println("=== Selector Fallback Chain ===")
    tree := behave.New(selector)
    status := tree.Tick()
    fmt.Printf("Final result: %s\n\n", status.String())
    
    // Selector where all options fail
    failingOption1 := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Option 1: Failed!")
            return behave.Failure
        },
    }
    
    failingOption2 := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Option 2: Failed!")
            return behave.Failure
        },
    }
    
    failingSelector := &behave.Selector{
        Children: []behave.Node{failingOption1, failingOption2},
    }
    
    fmt.Println("=== All Options Fail ===")
    tree2 := behave.New(failingSelector)
    status = tree2.Tick()
    fmt.Printf("Final result: %s\n\n", status.String())
    
    // Selector with running state
    runningTicks := 0
    runningOption := &behave.Action{
        Run: func() behave.Status {
            runningTicks++
            fmt.Printf("Running option tick %d: ", runningTicks)
            
            if runningTicks < 2 {
                fmt.Println("Still running...")
                return behave.Running
            }
            
            fmt.Println("Success!")
            return behave.Success
        },
    }
    
    notExecutedOption := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("This should not execute")
            return behave.Success
        },
    }
    
    runningSelector := &behave.Selector{
        Children: []behave.Node{runningOption, notExecutedOption},
    }
    
    fmt.Println("=== Selector with Running State ===")
    tree3 := behave.New(runningSelector)
    for tree3.Status() != behave.Success {
        status = tree3.Tick()
        fmt.Printf("Status: %s\n", status.String())
    }
}
```

## Parallel Node Example

The Parallel node runs all children simultaneously and succeeds when enough children succeed:

```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/rbrabson/behave"
)

func main() {
    // Create parallel tasks
    task1Ticks := 0
    task1 := &behave.Action{
        Run: func() behave.Status {
            task1Ticks++
            fmt.Printf("Task 1 tick %d: ", task1Ticks)
            
            if task1Ticks < 2 {
                fmt.Println("Running...")
                return behave.Running
            }
            
            fmt.Println("Success!")
            return behave.Success
        },
    }
    
    task2 := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Task 2: Quick success!")
            return behave.Success
        },
    }
    
    task3Ticks := 0
    task3 := &behave.Action{
        Run: func() behave.Status {
            task3Ticks++
            fmt.Printf("Task 3 tick %d: ", task3Ticks)
            
            if rand.Float32() < 0.5 { // 50% failure rate
                fmt.Println("Failed!")
                return behave.Failure
            }
            
            fmt.Println("Success!")
            return behave.Success
        },
    }
    
    // Parallel requiring 2 out of 3 successes
    parallel := &behave.Parallel{
        MinSuccessCount: 2,
        Children: []behave.Node{task1, task2, task3},
    }
    
    fmt.Println("=== Parallel Execution (2 out of 3 needed) ===")
    tree := behave.New(parallel)
    for tree.Status() == behave.Ready || tree.Status() == behave.Running {
        status := tree.Tick()
        fmt.Printf("Parallel status: %s\n", status.String())
        fmt.Println("---")
    }
    
    fmt.Printf("Final result: %s\n\n", tree.Status().String())
    
    // Parallel requiring all children to succeed
    task1Ticks = 0
    task3Ticks = 0
    
    allSucceedParallel := &behave.Parallel{
        MinSuccessCount: 3, // All must succeed
        Children: []behave.Node{task1, task2, task3},
    }
    
    fmt.Println("=== Parallel Execution (All must succeed) ===")
    tree2 := behave.New(allSucceedParallel)
    for i := 0; i < 5 && (tree2.Status() == behave.Ready || tree2.Status() == behave.Running); i++ {
        status := tree2.Tick()
        fmt.Printf("Parallel status: %s\n", status.String())
        if status != behave.Running {
            break
        }
        fmt.Println("---")
    }
    
    fmt.Printf("Final result: %s\n", tree2.Status().String())
}
```

## Composite Node Example

The Composite node combines a condition with another node, executing the child only if the condition succeeds:

```go
package main

import (
    "fmt"
    "github.com/rbrabson/behave"
)

func main() {
    // System state
    hasPermission := false
    resourceAvailable := true
    taskExecuted := false
    
    // Permission check condition
    permissionCheck := &behave.Condition{
        Check: func() behave.Status {
            fmt.Printf("Checking permission: %t\n", hasPermission)
            if hasPermission {
                fmt.Println("Permission granted!")
                return behave.Success
            }
            fmt.Println("Permission denied!")
            return behave.Failure
        },
    }
    
    // Protected action
    protectedAction := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Executing protected action...")
            taskExecuted = true
            return behave.Success
        },
    }
    
    // Composite: only execute action if permission check passes
    guardedExecution := &behave.Composite{
        Condition: permissionCheck,
        Child:     protectedAction,
    }
    
    fmt.Println("=== Composite without Permission ===")
    tree1 := behave.New(guardedExecution)
    status := tree1.Tick()
    fmt.Printf("Result: %s, Task executed: %t\n\n", status.String(), taskExecuted)
    
    // Grant permission and try again
    hasPermission = true
    taskExecuted = false
    
    fmt.Println("=== Composite with Permission ===")
    tree1.Reset()
    status = tree1.Tick()
    fmt.Printf("Result: %s, Task executed: %t\n\n", status.String(), taskExecuted)
    
    // More complex example with multiple conditions
    resourceCheck := &behave.Condition{
        Check: func() behave.Status {
            fmt.Printf("Resource available: %t\n", resourceAvailable)
            if resourceAvailable {
                return behave.Success
            }
            return behave.Failure
        },
    }
    
    // Multiple condition checks in sequence
    multiConditionCheck := &behave.Sequence{
        Children: []behave.Node{permissionCheck, resourceCheck},
    }
    
    criticalAction := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("Executing critical action...")
            return behave.Success
        },
    }
    
    // Composite with multiple conditions
    complexGuard := &behave.Composite{
        Condition: multiConditionCheck,
        Child:     criticalAction,
    }
    
    fmt.Println("=== Composite with Multiple Conditions ===")
    tree2 := behave.New(complexGuard)
    status = tree2.Tick()
    fmt.Printf("Result: %s\n\n", status.String())
    
    // Test with resource unavailable
    resourceAvailable = false
    
    fmt.Println("=== Composite with Resource Unavailable ===")
    tree2.Reset()
    status = tree2.Tick()
    fmt.Printf("Result: %s\n", status.String())
}
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

## WithTimeout Node Example

The WithTimeout decorator runs its child node for at most the specified duration. If the child completes (returns Success or Failure) before the duration expires, WithTimeout returns that status immediately. If the duration expires while the child is still running, WithTimeout returns Failure.

```go
package main

import (
    "fmt"
    "time"
    "github.com/rbrabson/behave"
)

func main() {
    ticks := 0
    timedAction := &behave.Action{
        Run: func() behave.Status {
            ticks++
            fmt.Printf("Tick %d\n", ticks)
            return behave.Success
        },
    }

    // Run the action for 1 second
    withTimeout := &behave.WithTimeout{
        Child:    timedAction,
        Duration: 1 * time.Second,
    }

    tree := behave.New(withTimeout)

    fmt.Println("=== WithTimeout Example ===")
    start := time.Now()
    for tree.Status() == behave.Ready || tree.Status() == behave.Running {
        status := tree.Tick()
        fmt.Printf("Status: %s\n", status.String())
        time.Sleep(200 * time.Millisecond)
    }
    elapsed := time.Since(start)
    fmt.Printf("Completed after %v and %d ticks!\n", elapsed, ticks)
}
```

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

## Forever Node Example

The Forever node is useful for running a child node indefinitely, regardless of its status. This is helpful for background tasks, polling, or infinite loops. It can also execute the root node indefinitely if the behavior tree should always be running.

```go
package main

import (
    "fmt"
    "time"
    "github.com/rbrabson/behave"
)

func main() {
    ticks := 0
    // Action that prints a message and always succeeds
    backgroundAction := &behave.Action{
        Run: func() behave.Status {
            ticks++
            fmt.Printf("Background tick %d\n", ticks)
            return behave.Success // Status is ignored by Forever
        },
    }

    forever := &behave.Forever{Child: backgroundAction}
    tree := behave.New(forever)

    // Simulate ticking the tree for a few cycles
    for i := 0; i < 5; i++ {
        status := tree.Tick()
        fmt.Printf("Tree status: %s\n", status.String())
        time.Sleep(500 * time.Millisecond)
    }

    fmt.Printf("Background ran for %d ticks!\n", ticks)
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

## Log Node Example

The Log node provides structured logging for debugging and monitoring behavior tree execution. It supports custom log levels or uses intelligent defaults:

- **Default Levels**: Info for Success, Warn for Failure, Debug for Running/Ready
- **Custom Levels**: Override defaults by setting the LogLevel field

```go
package main

import (
    "fmt"
    "log/slog"
    "math/rand"
    "os"
    "github.com/rbrabson/behave"
)

func main() {
    // Configure structured logging with JSON format
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    slog.SetDefault(logger)

    fmt.Println("=== Log Node Example ===")
    
    attempts := 0
    
    // Critical operation that might fail
    criticalOperation := &behave.Action{
        Run: func() behave.Status {
            attempts++
            fmt.Printf("Attempt %d: ", attempts)
            
            if rand.Float32() < 0.7 { // 70% chance of failure
                fmt.Println("Operation failed")
                return behave.Failure
            }
            
            fmt.Println("Operation succeeded")
            return behave.Success
        },
    }
    
    // Wrap with logging for monitoring
    loggedOperation := &behave.Log{
        Child:   criticalOperation,
        Message: "Critical system operation",
    }
    
    // Example with custom log level - always log at ERROR level for high priority
    highPriorityOperation := &behave.Action{
        Run: func() behave.Status {
            fmt.Println("High priority task executed")
            return behave.Success
        },
    }
    
    errorLevelLog := &behave.Log{
        Child:    highPriorityOperation,
        Message:  "High priority operation",
        LogLevel: func() *slog.Level { l := slog.LevelError; return &l }(), // Always ERROR level
    }
    
    // Retry with logging
    retryWithLogging := &behave.Retry{Child: loggedOperation}
    
    // Sequence: run both operations
    sequence := &behave.Sequence{
        Children: []behave.Node{retryWithLogging, errorLevelLog},
    }
    
    tree := behave.New(sequence)
    
    // Execute until success
    for tree.Status() != behave.Success {
        tree.Tick()
        if attempts > 10 { // Safety break
            break
        }
    }
    
    fmt.Printf("Operations completed after %d attempts\n", attempts)
    
    // Example 2: Complex behavior tree with multiple logged components
    fmt.Println("\n=== Complex Tree with Logging ===")
    
    // Database connection with logging (default levels)
    dbConnections := 0
    dbConnect := &behave.Log{
        Child: &behave.Action{
            Run: func() behave.Status {
                dbConnections++
                if dbConnections < 2 {
                    return behave.Failure
                }
                return behave.Success
            },
        },
        Message: "Database connection attempt",
        // Uses default levels: WARN for Failure, INFO for Success
    }
    
    // API call with custom INFO level (always log at INFO regardless of result)
    apiCalls := 0
    apiCall := &behave.Log{
        Child: &behave.Action{
            Run: func() behave.Status {
                apiCalls++
                return behave.Success
            },
        },
        Message:  "External API call",
        LogLevel: func() *slog.Level { l := slog.LevelInfo; return &l }(),
    }
    
    // Data processing with custom ERROR level (high importance)
    dataProcess := &behave.Log{
        Child: &behave.Action{
            Run: func() behave.Status {
                return behave.Success
            },
        },
        Message:  "Data processing step",
        LogLevel: func() *slog.Level { l := slog.LevelError; return &l }(),
    }
    
    // Sequence with logged steps
    pipeline := &behave.Sequence{
        Children: []behave.Node{
            &behave.Retry{Child: dbConnect},  // Retry DB connection
            apiCall,                          // Make API call
            dataProcess,                      // Process data
        },
    }
    
    pipelineTree := behave.New(pipeline)
    
    // Execute pipeline
    status := pipelineTree.Tick()
    fmt.Printf("Pipeline completed with status: %s\n", status.String())
    fmt.Printf("DB connections: %d, API calls: %d\n", dbConnections, apiCalls)
}
```

## Example: Using Struct Methods as Actions

You can use methods of your own struct as actions in a behavior tree. This allows you to encapsulate state and logic in a reusable way:

```go
package main

import (
    "fmt"
    "github.com/rbrabson/behave"
)

type ExampleActor struct {
    state int
}

// Always succeeds
func (a *ExampleActor) Succeed() behave.Status {
    a.state++
    return behave.Success
}

// Always fails
func (a *ExampleActor) Fail() behave.Status {
    a.state--
    return behave.Failure
}

// Succeeds if state is even, fails otherwise
func (a *ExampleActor) SucceedIfEven() behave.Status {
    if a.state%2 == 0 {
        return behave.Success
    }
    return behave.Failure
}

func main() {
    actor := &ExampleActor{state: 0}

    // Create Actions using the methods
    succeedNode := &behave.Action{Run: actor.Succeed}
    failNode := &behave.Action{Run: actor.Fail}
    evenNode := &behave.Action{Run: actor.SucceedIfEven}

    // Sequence: succeed, fail, even
    seq := &behave.Sequence{Children: []behave.Node{succeedNode, failNode, evenNode}}
    status := seq.Tick()
    fmt.Println("Sequence status:", status)

    // Selector: even, succeed
    sel := &behave.Selector{Children: []behave.Node{evenNode, succeedNode}}
    status2 := sel.Tick()
    fmt.Println("Selector status:", status2)
}
```

Output:

``` bash
Sequence status: Failure
Selector status: Success
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

// Log node for debugging and monitoring
loggedAction := &behave.Log{
    Child:   &behave.Action{Run: importantFunc},
    Message: "Executing critical operation",
    LogLevel: func() *slog.Level { l := slog.LevelError; return &l }(), // Custom error level
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
