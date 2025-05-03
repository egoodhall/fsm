# FSM - Finite State Machine Library

A Go library for building persistent, type-safe finite state machines. It provides a simple way to model complex workflows with guaranteed state transitions and persistence.

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

## Features

- Type-safe state transitions with generics
- Persistent state storage using SQLite
- Automatic state recovery on restart
- Thread-safe task processing
- Built-in error handling

## Development

```bash
# Install dependencies
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate SQL code
make generate

# Create new migration
make migration
```

## License

MIT
