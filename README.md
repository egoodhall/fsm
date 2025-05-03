# FSM - Finite State Machine Library

A Go library for building persistent, type-safe finite state machines. This library provides a simple way to model complex workflows with guaranteed state transitions and persistence.

## Features

- Type-safe state transitions with generics
- Persistent state storage using SQLite
- Automatic state recovery on restart
- Built-in error handling and state history
- Simple builder pattern API
- Thread-safe task processing

## Installation

```bash
go get github.com/egoodhall/fsm
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/egoodhall/fsm"
)

func main() {
    ctx := context.Background()
    
    // Create a new state machine
    machine, err := fsm.New[string, string]("example").
        // Define initial state transition
        InitialState(func(ctx context.Context, req fsm.FirstInput[string]) (fsm.Output, error) {
            if req.Input() != "start" {
                return fsm.Done()
            }
            return fsm.GoTo("processing", "data")
        }).
        // Add state transitions
        AddState("processing", func(ctx context.Context, req fsm.NthInput[string, string]) (fsm.Output, error) {
            return fsm.GoTo("completed", req.Previous()+" processed")
        }).
        Build(ctx, fsm.DB[string, string]("fsm.db"))

    if err != nil {
        panic(err)
    }

    // Start the machine
    submit, stop, err := machine.Start(ctx)
    if err != nil {
        panic(err)
    }
    defer stop()

    // Submit a task
    if err := submit(ctx, "task-1", "start"); err != nil {
        panic(err)
    }
}
```

## Concepts

### States
- `__initial__`: Starting state for new tasks
- `__done__`: Terminal state for completed tasks
- `__error__`: State for failed transitions
- Custom states: Define your own states for your workflow

### Transitions
- `FirstTransition`: Handles initial state transitions
- `NthTransition`: Handles transitions between custom states
- Each transition returns the next state and output data

### Persistence
- State transitions are stored in SQLite
- Tasks are automatically resumed on restart
- Full history of state transitions is maintained

## API Reference

### State Machine Builder
```go
fsm.New[IN, OUT](name string) InitialStateBuilder[IN, OUT]
```

### State Transitions
```go
// Initial state transition
func(ctx context.Context, req FirstInput[IN]) (Output, error)

// Subsequent state transitions
func(ctx context.Context, req NthInput[IN, OUT]) (Output, error)
```

### Output Helpers
```go
// Transition to a new state with data
GoTo[T any](state State, data T) (Output, error)

// Complete the task
Done() (Output, error)

// Mark task as failed
Error(err error) (Output, error)
```

## Development

### Prerequisites
```bash
# Install goose for migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Install sqlc for type-safe SQL
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Commands
```bash
# Generate SQL code
make generate

# Create new migration
make migration
```

## License

MIT
