# FSM - A Go Finite State Machine Generator

This project provides a DSL (Domain Specific Language) for defining finite state machines in Go.

## Installation

```bash
go get github.com/egoodhall/fsm
```

## Usage

### Define your FSM using the DSL

Create a file named `myfsm.fsm` with your FSM definition:

```
type MyInput;
type MyOutput;
type ContextA;
type ContextB;

// MyStateMachine handles transitions between states
fsm MyStateMachine[MyInput] {
  start A;
  state B[ContextA];
  state C[ContextB];
  state D;
  end END;

  transition A to B or C or END;
  transition B to D or END;
  transition C to D;
  transition D to END;
}
```

### Generate Go code from your FSM definition

```bash
make grammar  # Generate the parser
go run github.com/egoodhall/fsm/cmd/fsm-gen myfsm.fsm > myfsm.go
```

### Use the generated code in your application

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourmodule/myfsm"
)

func main() {
	ctx := context.Background()
	
	// Create the FSM
	machine, err := myfsm.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create FSM: %v", err)
	}
	
	// Start processing
	submit, stop, err := machine.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start FSM: %v", err)
	}
	defer stop()
	
	// Submit a task
	id, err := submit(ctx, 0, "some input")
	if err != nil {
		log.Fatalf("Failed to submit task: %v", err)
	}
	
	fmt.Printf("Task submitted with ID: %d\n", id)
}
```

## DSL Syntax

### Type Declarations

```
type TypeName;
```

### State Machine Definition

```
fsm MachineName[InputType] {
  // State declarations
  start StateName;
  state StateName[ContextType];
  state StateName;
  end StateName;

  // Transition declarations
  transition FromState to ToState1 or ToState2 or ToState3;
}
```

## License

[MIT License](LICENSE)
